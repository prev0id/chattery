package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type serverService interface {
	GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error)
	SearchServersNotJoined(ctx context.Context, userID domain.UserID, query string) ([]*domain.Server, error)

	CreateServer(ctx context.Context, name string, creatorUserID domain.UserID) (domain.ServerID, error)
	UpdateServer(ctx context.Context, serverID domain.ServerID, name string, userID domain.UserID) error
	DeleteServer(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error

	AddParticipant(ctx context.Context, participant *domain.ServerParticipant) error
	RemoveParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error

	CreateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) (domain.TopicID, error)
	UpdateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) error
	DeleteTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
}

type textTopicService interface {
	CreateMessage(ctx context.Context, message *domain.TopicMessage) error
	FirstPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error)
	NextPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error)
}

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type userCache interface {
	GetByID(id domain.UserID) (*domain.User, error)
	ListByID() map[domain.UserID]*domain.User
}

type Server struct {
	user      userService
	server    serverService
	textTopic textTopicService
	cache     userCache
}

func New(user userService, server serverService, textTopic textTopicService, cache userCache) *Server {
	return &Server{
		user:      user,
		server:    server,
		textTopic: textTopic,
		cache:     cache,
	}
}

func (*Server) Pattern() string {
	return "/v1/server"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Get("/list", s.GetServers)
		withAuthRouter.Get("/search", s.SearchServers)
		withAuthRouter.Post("/create", s.PostCreateServer)
		withAuthRouter.Post("/join", s.PostJoinServer)
		withAuthRouter.Post("/leave", s.PostLeaveServer)
		withAuthRouter.Post("/update", s.PostUpdateServer)
		withAuthRouter.Delete("/delete", s.DeleteServer)

		withAuthRouter.Post("/topic/create", s.PostTopicCreate)
		withAuthRouter.Post("/topic/update", s.PostTopicUpdate)
		withAuthRouter.Delete("/topic/delete", s.DeleteTopic)
		withAuthRouter.Post("/topic/message", s.PostTopicMessage)
		withAuthRouter.Post("/topic/messages", s.PostTopicMessages)
	})
}
