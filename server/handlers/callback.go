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

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Error().Msg("missing authorization code")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		log.Info().Msg("received oauth callback")

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

		fmt.Fprintf(w, "Authentication successful! Token saved. You can close this window.")
	}
}
