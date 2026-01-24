package chat_api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type chatService interface {
	UserChats(ctx context.Context, user domain.UserID) ([]*domain.Chat, error)
	JoinChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error
	LeaveChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error

	CreatePublicChat(ctx context.Context, user domain.UserID, name string) (domain.ChatID, error)
	CreatePrivateChat(ctx context.Context, users ...domain.UserID) (domain.ChatID, error)
	DeleteChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error

	SearchChats(ctx context.Context, query string) ([]*domain.Chat, error)

	ListMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.MessageCursor) ([]*domain.Message, *domain.MessageCursor, error)
}

type Server struct {
	user userService
	chat chatService
}

func New(user userService, chat chatService) *Server {
	return &Server{
		user: user,
		chat: chat,
	}
}

func (s *Server) Pattern() string {
	return "/v1/chat"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Post("/create/public", s.CreatePublic)
		withAuthRouter.Post("/create/private", s.CreatePrivate)
		withAuthRouter.Get("/search", s.Search)

		withAuthRouter.Get("/me/list", s.ListMy)
		withAuthRouter.Post("/me/join", s.Join)
		withAuthRouter.Delete("/me/leave", s.Leave)
	})
}
