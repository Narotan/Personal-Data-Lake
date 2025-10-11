package server

import (
	"DataLake/internal/middleware"
	"DataLake/server/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) routes() {
	s.mux.Handle("/callback", middleware.Logging(handlers.HandleCallback(s.cfg)))
	s.mux.Handle("/metrics", promhttp.Handler())
}
