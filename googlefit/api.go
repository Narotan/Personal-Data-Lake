package googlefit

import (
	"DataLake/auth"
	googlefitauth "DataLake/auth/googlefit"
	internal_db "DataLake/internal/db"
	googlefit_db "DataLake/internal/db/googlefit"
	"DataLake/internal/logger"
	"DataLake/internal/metrics"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	uuid "github.com/satori/go.uuid"
)

// FetchSummaries получает агрегированные данные за последние N дней
func FetchSummaries(days int) (*AggregatedDataResponse, error) {
	log := logger.Get()
	start := time.Now()
	metrics.GoogleFitFetchTotal.Inc()
	ctx := context.Background()

	storage := auth.NewFileTokenStorage("tokens.json")
	provider := googlefitauth.NewProviderFromEnv()
	tokenManager := auth.NewTokenManager(storage, provider)

	token, err := tokenManager.GetValidToken(ctx, "googlefit")
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to get valid token")
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}

	url := "https://www.googleapis.com/fitness/v1/users/me/dataset:aggregate"

	endTime := time.Now().UTC()
	startTime := endTime.AddDate(0, 0, -days)

	aggregateBy := map[string]interface{}{
		"aggregateBy": []map[string]string{
			{
				"dataTypeName": DataTypeStepCount,
			},
			{
				"dataTypeName": DataTypeDistance,
			},
		},
		"bucketByTime": map[string]int64{
			"durationMillis": 86400000, // 1 день
		},
		"startTimeMillis": startTime.UnixMilli(),
		"endTimeMillis":   endTime.UnixMilli(),
	}

	log.Info().
		Str("start_date", startTime.Format("2006-01-02")).
		Str("end_date", endTime.Format("2006-01-02")).
		Int("days", days).
		Msg("fetching google fit summaries")

	bodyBytes, err := json.Marshal(aggregateBy)
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to marshal request body")
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to create request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to execute request")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response", string(body)).
			Msg("unexpected status code")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response AggregatedDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	metrics.GoogleFitFetchDuration.Observe(time.Since(start).Seconds())
	log.Info().
		Int("buckets_fetched", len(response.Bucket)).
		Dur("duration", time.Since(start)).
		Msg("successfully fetched google fit summaries")

	return &response, nil
}

// SaveSummaries сохраняет агрегированные данные Google Fit в БД
func SaveSummaries(store *internal_db.Store, response *AggregatedDataResponse, userID uuid.UUID) error {
	log := logger.Get()
	start := time.Now()

	ctx := context.Background()

	var uuidBytes [16]byte
	copy(uuidBytes[:], userID.Bytes())

	stats, err := response.ExtractDailyStats()
	if err != nil {
		log.Error().Err(err).Msg("failed to extract daily stats")
		return fmt.Errorf("failed to extract daily stats: %w", err)
	}

	log.Info().
		Int("days_count", len(stats)).
		Str("user_id", userID.String()).
		Msg("saving google fit daily stats")

	for _, stat := range stats {
		err := store.ExecTxGoogleFit(ctx, func(q *googlefit_db.Queries) error {
			_, err := q.UpsertDailyStat(ctx, googlefit_db.UpsertDailyStatParams{
				UserID:   pgtype.UUID{Bytes: uuidBytes, Valid: true},
				Date:     pgtype.Date{Time: stat.Date, Valid: true},
				Steps:    pgtype.Int4{Int32: int32(stat.Steps), Valid: true},
				Distance: pgtype.Float8{Float64: stat.Distance, Valid: true},
			})
			if err != nil {
				metrics.DatabaseOperationsTotal.WithLabelValues("upsert_daily_stat", "error").Inc()
				log.Error().
					Err(err).
					Str("date", stat.Date.Format("2006-01-02")).
					Msg("failed to upsert daily stat")
				return fmt.Errorf("failed to upsert daily stat for %s: %w", stat.Date.Format("2006-01-02"), err)
			}
			metrics.DatabaseOperationsTotal.WithLabelValues("upsert_daily_stat", "success").Inc()

			log.Debug().
				Str("date", stat.Date.Format("2006-01-02")).
				Int("steps", stat.Steps).
				Float64("distance", stat.Distance).
				Msg("saved daily stat")

			return nil
		})

		if err != nil {
			return err
		}
	}

	metrics.DatabaseOperationDuration.WithLabelValues("save_googlefit_summaries").Observe(time.Since(start).Seconds())
	log.Info().
		Int("days_saved", len(stats)).
		Dur("duration", time.Since(start)).
		Msg("successfully saved all google fit summaries")

	return nil
}
