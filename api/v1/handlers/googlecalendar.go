package handlers_api_v1

import (
	models_api_v1 "DataLake/api/v1/models"
	internal_db "DataLake/internal/db"
	googlecalendar_db "DataLake/internal/db/googlecalendar"
	"DataLake/internal/middleware"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type GoogleCalendarHandler struct {
	store  *internal_db.Store
	logger *zerolog.Logger
}

func NewGoogleCalendarHandler(store *internal_db.Store, logger *zerolog.Logger) *GoogleCalendarHandler {
	return &GoogleCalendarHandler{
		store:  store,
		logger: logger,
	}
}

// GetEvents обрабатывает GET /api/v1/googlecalendar/events.

func (h *GoogleCalendarHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -6)

	if val := r.URL.Query().Get("start_date"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			startDate = t
		} else {
			http.Error(w, `{"error": "Invalid start_date format. Use YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
	}
	if val := r.URL.Query().Get("end_date"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			endDate = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // Включаем весь день
		} else {
			http.Error(w, `{"error": "Invalid end_date format. Use YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
	}

	userIDStr, ok := middleware.GetUserID(r.Context())
	if !ok || userIDStr == "" {
		h.logger.Error().Msg("Failed to get user ID from context")
		http.Error(w, `{"error": "Unauthorized: User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid user ID format")
		http.Error(w, `{"error": "Internal Server Error: Invalid user ID format"}`, http.StatusInternalServerError)
		return
	}

	var userIDBytes [16]byte
	copy(userIDBytes[:], userID.Bytes())

	dbResult, err := h.store.GoogleCalendar.GetCalendarEventsByDateRange(r.Context(), googlecalendar_db.GetCalendarEventsByDateRangeParams{
		UserID:      pgtype.UUID{Bytes: userIDBytes, Valid: true},
		StartTime:   pgtype.Timestamptz{Time: startDate, Valid: true},
		StartTime_2: pgtype.Timestamptz{Time: endDate, Valid: true},
	})

	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get google calendar events from DB")
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := make([]models_api_v1.CalendarEvent, 0, len(dbResult))
	for _, row := range dbResult {
		response = append(response, models_api_v1.CalendarEvent{
			ID:          row.EventID,
			Summary:     row.Summary.String,
			Description: row.Description.String,
			StartTime:   row.StartTime.Time.Format(time.RFC3339),
			EndTime:     row.EndTime.Time.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
