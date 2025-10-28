package server

import (
	internal_db "DataLake/internal/db"
	"DataLake/internal/logger"
	"net/http"
)

type Server struct {
	store *internal_db.Store
	mux   *http.ServeMux
}

func NewServer(store *internal_db.Store) *Server {
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

func (s *Server) Store() *internal_db.Store {
	return s.store
}
