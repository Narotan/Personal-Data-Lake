package server

import (
	"DataLake/internal/middleware"
	"DataLake/server/handlers"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) routes(apiRouter http.Handler) {
	// Health check (с CORS)
	s.mux.Handle("/health", middleware.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})))

	// Setup page (OAuth instructions)
	s.mux.Handle("/", middleware.CORS(middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/setup.html")
	}))))

	// WakaTime OAuth
	s.mux.Handle("/callback", middleware.CORS(middleware.Logging(handlers.HandleCallback())))

	// Google Fit OAuth
	s.mux.Handle("/auth/googlefit", middleware.CORS(middleware.Logging(handlers.HandleGoogleFitAuth())))
	s.mux.Handle("/oauth2callback", middleware.CORS(middleware.Logging(handlers.HandleGoogleFitCallback())))

	// Google Calendar OAuth
	s.mux.Handle("/auth/googlecalendar", middleware.CORS(middleware.Logging(handlers.HandleGoogleCalendarAuth())))
	s.mux.Handle("/oauth2callback/calendar", middleware.CORS(middleware.Logging(handlers.HandleGoogleCalendarCallback())))

	// API v1 (с CORS и Rate Limiting)
	s.mux.Handle("/api/v1/auth/status", middleware.RateLimit(middleware.CORS(middleware.Logging(handlers.HandleAuthStatus()))))
	s.mux.Handle("/api/v1/", middleware.RateLimit(middleware.CORS(http.StripPrefix("/api/v1", apiRouter))))

	// Metrics
	s.mux.Handle("/metrics", promhttp.Handler())
}
