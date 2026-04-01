package server

import (
	"context"

	"chattery/internal/config"
	"chattery/internal/domain"
)

type db interface {
	GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error)

	GetServer(ctx context.Context, serverID domain.ServerID) (*domain.Server, error)
	CreateServer(ctx context.Context, name string) (domain.ServerID, error)
	UpdateServer(ctx context.Context, server *domain.Server) error
	DeleteServer(ctx context.Context, serverID domain.ServerID) error

	GetServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) (*domain.ServerParticipant, error)
	CreateServerParticipant(ctx context.Context, participant *domain.ServerParticipant) error
	DeleteServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error

	GetTopic(ctx context.Context, topicID domain.TopicID) (*domain.Topic, error)
	CreateTopic(ctx context.Context, topic *domain.Topic) (domain.TopicID, error)
	UpdateTopic(ctx context.Context, topic *domain.Topic) error
	DeleteTopic(ctx context.Context, topicID domain.TopicID) error

	CreateMessage(ctx context.Context, message *domain.TopicMessage) error
	FirstPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error)
	NextPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error)
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db          db
	transaction txManager
	limit       int
}

func New(dbAdapter db, transaction txManager, cfg *config.Config) *Service {
	return &Service{
		db:          dbAdapter,
		transaction: transaction,
		limit:       cfg.Chat.MessagesLimit,
	}
}

func (s *Service) getNextCursor(topicID domain.TopicID, messages []*domain.TopicMessage) *domain.TopicCursor {
	if len(messages) != s.limit {
		return nil
	}

	lastMessage := messages[len(messages)-1]
	return &domain.TopicCursor{
		ChatID:    topicID,
		MessageID: lastMessage.ID,
		Timestamp: lastMessage.CreatedAt,
	}
}
