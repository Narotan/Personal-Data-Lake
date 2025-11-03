package handlers_api_v1

import (
	internal_db "DataLake/internal/db"
	activitywatch_db "DataLake/internal/db/activitywatch"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type ActivityWatchHandler struct {
	store  *internal_db.Store
	logger *zerolog.Logger
}

func NewActivityWatchHandler(store *internal_db.Store, logger *zerolog.Logger) *ActivityWatchHandler {
	return &ActivityWatchHandler{
		store:  store,
		logger: logger,
	}
}

type ActivityEventRequest struct {
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration"`
	App       string    `json:"app"`
	Title     string    `json:"title"`
	BucketID  string    `json:"bucket_id"`
}

// HandleEvents обрабатывает POST /api/v1/activitywatch/events
func (h *ActivityWatchHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var events []ActivityEventRequest
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		h.logger.Error().Err(err).Msg("Failed to decode events")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(events) == 0 {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "No events to insert"})
		return
	}

	params := make([]activitywatch_db.BulkInsertEventsParams, len(events))
	for i, event := range events {
		params[i] = activitywatch_db.BulkInsertEventsParams{
			Timestamp: pgtype.Timestamptz{Time: event.Timestamp, Valid: true},
			Duration:  event.Duration,
			App:       event.App,
			Title:     pgtype.Text{String: event.Title, Valid: event.Title != ""},
			BucketID:  event.BucketID,
		}
	}

	ctx := context.Background()
	count, err := h.store.ActivityWatch.BulkInsertEvents(ctx, params)
	if err != nil {
		h.logger.Error().Err(err).Int("count", len(events)).Msg("Failed to insert events")
		http.Error(w, "Failed to save events", http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int64("count", count).Msg("Inserted activity events")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Events saved successfully",
		"inserted": count,
	})
}

// GetStats обрабатывает GET /api/v1/activitywatch/stats
func (h *ActivityWatchHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	startStr := query.Get("start")
	endStr := query.Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start time", http.StatusBadRequest)
			return
		}
	} else {
		start = time.Now().Add(-24 * time.Hour)
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end time", http.StatusBadRequest)
			return
		}
	} else {
		end = time.Now()
	}

	ctx := context.Background()
	params := activitywatch_db.GetAppStatsParams{
		Timestamp:   pgtype.Timestamptz{Time: start, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: end, Valid: true},
	}

	stats, err := h.store.ActivityWatch.GetAppStats(ctx, params)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get stats")
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}
