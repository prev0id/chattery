package server_api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type serverService interface {
	GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error)
	CreateServer(ctx context.Context, name string, creatorUserID domain.UserID) (domain.ServerID, error)
	UpdateServer(ctx context.Context, serverID domain.ServerID, name string) error
	DeleteServer(ctx context.Context, serverID domain.ServerID) error
	AddServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID, role domain.ServerRole) error
	RemoveServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error
	CreateTopic(ctx context.Context, serverID domain.ServerID, name string, topicType domain.TopicType) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, topicID domain.TopicID, name string) error
	DeleteTopic(ctx context.Context, topicID domain.TopicID) error
	CreateTopicMessage(ctx context.Context, topicID domain.TopicID, userID domain.UserID, text string) (domain.TopicMessageID, error)
	FirstPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, *domain.TopicCursor, error)
	NextPagesOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, *domain.TopicCursor, error)
}

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	user   userService
	server serverService
}

func New(user userService, server serverService) *Server {
	return &Server{
		user:   user,
		server: server,
	}
}

func (s *Server) Pattern() string {
	return "/v1/server"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Get("/list", s.List)
		withAuthRouter.Post("/create", s.Create)
		withAuthRouter.Post("/join", s.Join)
		withAuthRouter.Post("/leave", s.Leave)
		withAuthRouter.Post("/update", s.Update)
		withAuthRouter.Delete("/delete", s.Delete)

		withAuthRouter.Post("/topic/create", s.CreateTopic)
		withAuthRouter.Post("/topic/update", s.UpdateTopic)
		withAuthRouter.Delete("/topic/delete", s.DeleteTopic)
		withAuthRouter.Post("/topic/message", s.CreateTopicMessage)
		withAuthRouter.Get("/topic/messages", s.ListTopicMessages)
	})
}
