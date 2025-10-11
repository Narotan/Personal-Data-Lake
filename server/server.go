package server

import (
	"DataLake/auth"
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/internal/logger"
	"net/http"
)

type Server struct {
	cfg   auth.Config
	store *wakatime_db.Store
	mux   *http.ServeMux
}

func NewServer(cfg auth.Config, store *wakatime_db.Store) *Server {
	s := &Server{
		cfg:   cfg,
		store: store,
		mux:   http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) Run() error {
	log := logger.Get()
	log.Info().Msg("starting server on :8080")
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) Cfg() auth.Config {
	return s.cfg
}

func (s *Server) Queries() *wakatime_db.Queries {
	return s.store.Queries
}

func (s *Server) Store() *wakatime_db.Store {
	return s.store
}
