package handlers

import (
	"DataLake/auth"
	"DataLake/internal/logger"
	"encoding/json"
	"net/http"
)

type AuthStatusResponse struct {
	WakaTime       bool `json:"wakatime"`
	GoogleFit      bool `json:"googlefit"`
	GoogleCalendar bool `json:"googlecalendar"`
}

func HandleAuthStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		storage, err := auth.NewFileTokenStorageFromEnv("tokens.json")
		if err != nil {
			log.Error().Err(err).Msg("failed to initialize token storage")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		status := AuthStatusResponse{}

		if _, err := storage.LoadToken("wakatime"); err == nil {
			status.WakaTime = true
		}

		if _, err := storage.LoadToken("googlefit"); err == nil {
			status.GoogleFit = true
		}

		if _, err := storage.LoadToken("googlecalendar"); err == nil {
			status.GoogleCalendar = true
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			log.Error().Err(err).Msg("failed to encode auth status response")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
