package googlefit

import (
	"DataLake/auth"
	"DataLake/internal/logger"
	"DataLake/internal/metrics"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func FetchSummaries() string {
	log := logger.Get()
	start := time.Now()
	metrics.GoogleFitFetchTotal.Inc()

	storage := auth.NewFileTokenStorage("tokens.json")
	token, err := storage.LoadToken("googlefit")
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to load tokens")
		return ""
	}

	url := "https://www.googleapis.com/fitness/v1/users/me/dataset:aggregate"

	aggregateBy := map[string]interface{}{
		"aggregateBy": []map[string]string{
			{
				"dataTypeName": "com.google.step_count.delta",
				"dataSourceId": "derived:com.google.step_count.delta:com.google.android.gms:estimated_steps",
			},
		},
		"bucketByTime": map[string]int64{
			"durationMillis": 86400000,
		},
		"startTimeMillis": start.UnixMilli() - 7*86400000,
		"endTimeMillis":   start.UnixMilli(),
	}

	bodyBytes, err := json.Marshal(aggregateBy)
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to marshal request body")
		return ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to create request")
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to execute request")
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.GoogleFitFetchErrors.Inc()
		log.Error().Err(err).Msg("failed to read response body")
		return ""
	}

	log.Info().RawJSON("response", body).Msg("Google Fit raw response")

	return string(body)
}
