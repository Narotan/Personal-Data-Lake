package server

import (
	wakatime_db "DataLake/internal/db/wakatime"
	"DataLake/internal/logger"
	"net/http"
)

type Server struct {
	store *wakatime_db.Store
	mux   *http.ServeMux
}

func NewServer(store *wakatime_db.Store) *Server {
	s := &Server{
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

func (s *Server) Queries() *wakatime_db.Queries {
	return s.store.Queries
}

func (s *Server) Store() *wakatime_db.Store {
	return s.store
}
