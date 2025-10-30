package handlers

import (
	"DataLake/auth"
	googlecalendarauth "DataLake/auth/googlecalendar"
	"DataLake/internal/logger"
	"fmt"
	"net/http"
	"os"
)

func HandleGoogleCalendarCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		code := r.URL.Query().Get("code")
		if code == "" {
			log.Error().Msg("missing authorization code")
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		log.Info().Msg("received google calendar oauth callback")

		googlecalendarProvider := googlecalendarauth.NewProvider(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALENDAR_REDIRECT_URI"),
		)

		token, err := googlecalendarProvider.ExchangeToken(r.Context(), code)
		if err != nil {
			log.Error().Err(err).Msg("failed to exchange code for token")
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		storage := auth.NewFileTokenStorage("tokens.json")
		if err := storage.SaveToken("googlecalendar", token); err != nil {
			log.Error().Err(err).Msg("failed to save token")
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		log.Info().
			Str("access_token_prefix", token.AccessToken[:10]+"...").
			Str("expires_at", token.ExpiresAt).
			Msg("successfully saved google calendar token")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>✅ Google Calendar Authorization Successful!</h1>
				<p>Access token saved successfully.</p>
				<p>You can now close this window and test the API at <a href="/googlecalendar/test">/googlecalendar/test</a></p>
			</body>
			</html>
		`)
	}
}

func HandleGoogleCalendarAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		googlecalendarProvider := googlecalendarauth.NewProvider(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALENDAR_REDIRECT_URI"),
		)

		authURL := googlecalendarProvider.GetAuthURL("state")
		log.Info().Str("auth_url", authURL).Msg("redirecting to google calendar authorization")

		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}
