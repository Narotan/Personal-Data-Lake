package handlers

import (
	"DataLake/auth"
	wakatimeauth "DataLake/auth/wakatime"
	"DataLake/internal/logger"
	"fmt"
	"net/http"
	"os"
)

func HandleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		// Логируем весь URL для отладки
		log.Info().
			Str("full_url", r.URL.String()).
			Str("query", r.URL.RawQuery).
			Msg("received callback request")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errorParam := r.URL.Query().Get("error")
		errorDescription := r.URL.Query().Get("error_description")

		// Если WakaTime вернул ошибку
		if errorParam != "" {
			log.Error().
				Str("error", errorParam).
				Str("error_description", errorDescription).
				Msg("oauth error from wakatime")
			http.Error(w, fmt.Sprintf("OAuth Error: %s - %s", errorParam, errorDescription), http.StatusBadRequest)
			return
		}

		if code == "" {
			log.Error().Msg("missing authorization code")
			http.Error(w, "Missing code parameter. Please go through the authorization flow: check logs for auth_url", http.StatusBadRequest)
			return
		}

		log.Info().
			Str("code_preview", code[:min(10, len(code))]+"...").
			Str("state", state).
			Msg("received oauth callback with code")

		wakatimeProvider := wakatimeauth.NewProvider(
			os.Getenv("CLIENT_ID"),
			os.Getenv("CLIENT_SECRET"),
			os.Getenv("REDIRECT_URI"),
		)

		token, err := wakatimeProvider.ExchangeToken(r.Context(), code)
		if err != nil {
			log.Error().Err(err).Msg("failed to exchange code for token")
			http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		storage := auth.NewFileTokenStorage("tokens.json")
		err = storage.SaveToken("wakatime", token)

		if err != nil {
			log.Error().Err(err).Msg("failed to save tokens")
			http.Error(w, "Failed to save tokens", http.StatusInternalServerError)
			return
		}

		log.Info().Str("uid", token.UID).Msg("oauth flow completed successfully")

		http.Redirect(w, r, "http://localhost:8000/?auth_success=true", http.StatusTemporaryRedirect)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
