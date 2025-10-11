package handlers

import (
	"DataLake/auth"
	"DataLake/internal/logger"
	"fmt"
	"net/http"
)

func HandleCallback(cfg auth.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Error().Msg("missing authorization code")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		log.Info().Msg("received oauth callback")

		token, err := auth.ExchangeCodeForToken(cfg, code)
		if err != nil {
			log.Error().Err(err).Msg("failed to exchange code for token")
			http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = auth.SaveTokens(auth.Tokens{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt,
		})

		if err != nil {
			log.Error().Err(err).Msg("failed to save tokens")
			http.Error(w, "Failed to save tokens", http.StatusInternalServerError)
			return
		}

		log.Info().Str("uid", token.UID).Msg("oauth flow completed successfully")

		fmt.Fprintf(w, "Authentication successful! Token saved. You can close this window.")
	}
}
