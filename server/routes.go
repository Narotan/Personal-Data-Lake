package server

import (
	"DataLake/internal/middleware"
	"DataLake/server/handlers"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) routes(apiRouter http.Handler) {
	// Health check
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Setup page (OAuth instructions)
	s.mux.Handle("/", middleware.Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/setup.html")
	})))

	// WakaTime OAuth
	s.mux.Handle("/callback", middleware.Logging(handlers.HandleCallback()))

	// Google Fit OAuth
	s.mux.Handle("/auth/googlefit", middleware.Logging(handlers.HandleGoogleFitAuth()))
	s.mux.Handle("/oauth2callback", middleware.Logging(handlers.HandleGoogleFitCallback()))

	// Google Calendar OAuth
	s.mux.Handle("/auth/googlecalendar", middleware.Logging(handlers.HandleGoogleCalendarAuth()))
	s.mux.Handle("/oauth2callback/calendar", middleware.Logging(handlers.HandleGoogleCalendarCallback()))

	// API v1
	s.mux.Handle("/api/v1/auth/status", middleware.Logging(handlers.HandleAuthStatus()))
	s.mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	// Metrics
	s.mux.Handle("/metrics", promhttp.Handler())
}
