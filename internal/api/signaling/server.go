package signalingapi

import (
	"context"

	"github.com/go-chi/chi/v5"

	"chattery/internal/service/subscription"
)

type chatService interface {
	Register(ctx context.Context, subscriber *subscription.Subscriber)
}

type Server struct {
	chat chatService
}

func New() *Server {
	return &Server{}
}

func (s *Server) Pattern() string {
	return "/v1/signaling"
}

func (s *Server) Route(router chi.Router) {
	router.Get("/ws", s.WebsocketEntrypoint)
}
