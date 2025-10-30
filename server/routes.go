package server

import (
	"DataLake/internal/middleware"
	"DataLake/server/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) routes() {
	// WakaTime OAuth
	s.mux.Handle("/callback", middleware.Logging(handlers.HandleCallback()))

	// Google Fit OAuth
	s.mux.Handle("/auth/googlefit", middleware.Logging(handlers.HandleGoogleFitAuth()))
	s.mux.Handle("/oauth2callback", middleware.Logging(handlers.HandleGoogleFitCallback()))

	// Google Calendar OAuth
	s.mux.Handle("/auth/googlecalendar", middleware.Logging(handlers.HandleGoogleCalendarAuth()))
	s.mux.Handle("/oauth2callback/calendar", middleware.Logging(handlers.HandleGoogleCalendarCallback()))

	// Metrics
	s.mux.Handle("/metrics", promhttp.Handler())
}
