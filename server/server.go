package server

import (
	"DataLake/auth"
	wakatime_db "DataLake/internal/db/wakatime"
	"log"
	"net/http"
)

type Server struct {
	cfg     auth.Config
	queries *wakatime_db.Queries
	logger  *log.Logger
	mux     *http.ServeMux
}

func NewServer(cfg auth.Config, queries *wakatime_db.Queries, logger *log.Logger) *Server {
	s := &Server{
		cfg:     cfg,
		queries: queries,
		logger:  logger,
		mux:     http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) Run() error {
	s.logger.Println("Starting server on :8080")
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) Cfg() auth.Config {
	return s.cfg
}

func (s *Server) Queries() *wakatime_db.Queries {
	return s.queries
}

func (s *Server) Logger() *log.Logger {
	return s.logger
}
