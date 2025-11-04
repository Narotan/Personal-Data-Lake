package handlers_api_v1

import (
	models_api_v1 "DataLake/api/v1/models"
	internal_db "DataLake/internal/db"
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/internal/middleware"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type WakatimeHandler struct {
	store  *internal_db.Store
	logger *zerolog.Logger
}

func NewWakatimeHandler(store *internal_db.Store, logger *zerolog.Logger) *WakatimeHandler {
	return &WakatimeHandler{
		store:  store,
		logger: logger,
	}
}

// GetStats обрабатывает GET /api/v1/wakatime/stats.
// возвращает ежедневную статистику по времени проведенному в коде за указанный диапазон дат.
func (h *WakatimeHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -6)

	if val := r.URL.Query().Get("start_date"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			startDate = t
		} else {
			h.logger.Error().Err(err).Str("start_date", val).Msg("invalid start_date format")
			http.Error(w, `{"error": "Invalid start_date format. Use YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
	}

	if val := r.URL.Query().Get("end_date"); val != "" {
		if t, err := time.Parse("2006-01-02", val); err == nil {
			endDate = t
		} else {
			h.logger.Error().Err(err).Str("end_date", val).Msg("invalid end_date format")
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

	dbResult, err := h.store.WakaTime.GetWakatimeStatsByDateRange(r.Context(), wakatime_db.GetWakatimeStatsByDateRangeParams{
		UserID: pgtype.UUID{Bytes: userIDBytes, Valid: true},
		Date:   pgtype.Date{Time: startDate, Valid: true},
		Date_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get wakatime stats from DB")
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	dailyStatsMap := make(map[string]*models_api_v1.DailyStat)
	for _, row := range dbResult {
		if !row.Date.Valid {
			continue
		}
		dayDate := row.Date.Time.Format("2006-01-02")
		if _, ok := dailyStatsMap[dayDate]; !ok {
			dailyStatsMap[dayDate] = &models_api_v1.DailyStat{
				Date:         dayDate,
				TotalSeconds: row.TotalSeconds,
				Text:         row.DayText.String,
				Projects:     make([]models_api_v1.ProjectStat, 0),
			}
		}
		if row.ProjectName.Valid {
			dailyStatsMap[dayDate].Projects = append(dailyStatsMap[dayDate].Projects, models_api_v1.ProjectStat{
				Name:         row.ProjectName.String,
				TotalSeconds: row.ProjectSeconds.Float64,
			})
		}
	}

	response := make([]models_api_v1.DailyStat, 0, len(dailyStatsMap))
	for _, stat := range dailyStatsMap {
		response = append(response, *stat)
	}

	// Сортировка по дате (от новых к старым)
	sort.Slice(response, func(i, j int) bool {
		return response[i].Date > response[j].Date
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error().Err(err).Msg("failed to encode response")
	}
}
