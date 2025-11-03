package handlers_api_v1

import (
	models_api_v1 "DataLake/api/v1/models"
	internal_db "DataLake/internal/db"
	googlefit_db "DataLake/internal/db/googlefit"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type GoogleFitHandler struct {
	store  *internal_db.Store
	logger *zerolog.Logger
}

func NewGoogleFitHandler(store *internal_db.Store, logger *zerolog.Logger) *GoogleFitHandler {
	return &GoogleFitHandler{
		store:  store,
		logger: logger,
	}
}

// GetStats обрабатывает GET /api/v1/googlefit/stats.
// Возвращает ежедневную статистику (количество шагов, расстояние пройденное) за указанный диапазон дат.
func (h *GoogleFitHandler) GetStats(w http.ResponseWriter, r *http.Request) {
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
			endDate = t
		} else {
			http.Error(w, `{"error": "Invalid end_date format. Use YYYY-MM-DD"}`, http.StatusBadRequest)
			return
		}
	}

	userID := uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))
	var uidArr [16]byte
	copy(uidArr[:], userID.Bytes())
	dbResult, err := h.store.GoogleFit.GetGoogleFitDailyStatsByDateRange(r.Context(), googlefit_db.GetGoogleFitDailyStatsByDateRangeParams{
		UserID: pgtype.UUID{Bytes: uidArr, Valid: true},
		Date:   pgtype.Date{Time: startDate, Valid: true},
		Date_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get google fit stats from DB")
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := make([]models_api_v1.DailyFitStat, 0, len(dbResult))
	for _, row := range dbResult {
		response = append(response, models_api_v1.DailyFitStat{
			Date:     row.Date.Time.Format("2006-01-02"),
			Steps:    int(row.Steps.Int32),
			Distance: row.Distance.Float64,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
