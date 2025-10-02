package server

import (
	"DataLake/server/handlers"
)

func (s *Server) routes() {
	s.mux.Handle("/callback", handlers.HandleCallback(s.cfg, s.logger))
}
