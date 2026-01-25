package signaling_api

import (
	"context"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type chatService interface {
	Register(sub domain.Subscriber)
	Unregister(sub domain.Subscriber)

	PostMessage(ctx context.Context, message *domain.Message) error
	StartListeningToChat(ctx context.Context, sub domain.Subscriber, chat domain.ChatID)
	StopListeningToChat(sub domain.Subscriber)
}

type Server struct {
	chat chatService
}

func New(chat chatService) *Server {
	return &Server{
		chat: chat,
	}
}

func (s *Server) Pattern() string {
	return "/v1/signaling"
}

func (s *Server) Route(router chi.Router) {
	router.Get("/ws", s.WebsocketEntrypoint)
}
