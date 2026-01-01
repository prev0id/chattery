package chatapi

import (
	"github.com/go-chi/chi/v5"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) Pattern() string {
	return "/v1/chat"
}

func (s *Server) Route(router chi.Router) {
	router.Get("/ws", s.WebsocketEntrypoint)
}
