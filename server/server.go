package server

import (
	internal_db "DataLake/internal/db"
	"DataLake/internal/logger"
	"net/http"

	"github.com/rs/zerolog"
)

type Server struct {
	store  *internal_db.Store
	mux    *http.ServeMux
	logger zerolog.Logger
}

func NewServer(store *internal_db.Store) *Server {
	log := logger.Get()
	s := &Server{
		store:  store,
		mux:    http.NewServeMux(),
		logger: log,
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
