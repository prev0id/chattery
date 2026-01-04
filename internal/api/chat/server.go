package chatapi

import (
	"chattery/internal/domain"
	"chattery/internal/service/signaling"
	"context"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
)

type signalingService interface {
	Subscribe(ctx context.Context, user domain.Username, ws *websocket.Conn) (context.Context, *signaling.Subscriber)
	Unsubscribe(sub *signaling.Subscriber)
}

type chatService interface {
	Register(ctx context.Context, subscriber *signaling.Subscriber)
}

type Server struct {
	signaling signalingService
	chat      chatService
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
