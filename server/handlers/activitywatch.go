package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	activitywatch_db "DataLake/internal/db/activitywatch"

	"github.com/jackc/pgx/v5/pgtype"
)

type ActivityWatchHandler struct {
	queries *activitywatch_db.Queries
	logger  *slog.Logger
}

func NewActivityWatchHandler(queries *activitywatch_db.Queries, logger *slog.Logger) *ActivityWatchHandler {
	return &ActivityWatchHandler{
		queries: queries,
		logger:  logger,
	}
}

type ActivityEventPayload struct {
	ID        int64                  `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  float64                `json:"duration"`
	Data      map[string]interface{} `json:"data"`
}

type ActivityWatchRequest struct {
	BucketID string                 `json:"bucket_id"`
	Events   []ActivityEventPayload `json:"events"`
}

func (h *ActivityWatchHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ActivityWatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if len(req.Events) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "inserted": "0"})
		return
	}

	ctx := context.Background()

	for _, event := range req.Events {
		app := "unknown"
		title := "unknown"

		if appVal, ok := event.Data["app"].(string); ok {
			app = appVal
		}
		if titleVal, ok := event.Data["title"].(string); ok {
			title = titleVal
		}

		params := activitywatch_db.InsertActivityEventParams{
			ID:       event.ID,
			Duration: event.Duration,
			App:      app,
			BucketID: req.BucketID,
		}

		params.Timestamp = pgtype.Timestamptz{
			Time:  event.Timestamp,
			Valid: true,
		}

		params.Title = pgtype.Text{
			String: title,
			Valid:  title != "",
		}

		if err := h.queries.InsertActivityEvent(ctx, params); err != nil {
			h.logger.Error("failed to insert event", "error", err, "event_id", event.ID)
			continue
		}
	}

	h.logger.Info("inserted activity events", "count", len(req.Events), "bucket", req.BucketID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"inserted": len(req.Events),
	})
}
