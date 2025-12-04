package handlers

import (
	"DataLake/auth"
	googlefitauth "DataLake/auth/googlefit"
	"DataLake/internal/logger"
	"net/http"
	"os"
)

func HandleGoogleFitCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Error().Msg("missing authorization code")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		log.Info().Msg("received google fit oauth callback")

		googlefitProvider := googlefitauth.NewProvider(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URI"),
		)

		token, err := googlefitProvider.ExchangeToken(r.Context(), code)
		if err != nil {
			log.Error().Err(err).Msg("failed to exchange code for token")
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		storage := auth.NewFileTokenStorage("tokens.json")
		if err := storage.SaveToken("googlefit", token); err != nil {
			log.Error().Err(err).Msg("failed to save token")
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		log.Info().
			Str("access_token_prefix", token.AccessToken[:10]+"...").
			Str("expires_at", token.ExpiresAt).
			Msg("successfully saved google fit token")

		http.Redirect(w, r, "http://localhost:8000/?auth_success=true", http.StatusTemporaryRedirect)
	}
}

func HandleGoogleFitAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		googlefitProvider := googlefitauth.NewProvider(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_REDIRECT_URI"),
		)

		authURL := googlefitProvider.GetAuthURL("state")
		log.Info().Str("auth_url", authURL).Msg("redirecting to google fit authorization")

		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}
