package server

import (
	"context"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
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
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type redis interface {
	PublishToUser(ctx context.Context, userID domain.UserID, message *domain.UserMessage) error
}

type Service struct {
	db          db
	transaction txManager
	redis       redis
	limit       int
}

func New(dbAdapter db, transaction txManager, redisAdapter redis, cfg *config.Config) *Service {
	return &Service{
		db:          dbAdapter,
		transaction: transaction,
		redis:       redisAdapter,
		limit:       cfg.Chat.MessagesLimit,
	}
}

func (s *Service) getTopic(ctx context.Context, topicID domain.TopicID) (*domain.Topic, error) {
	topic, err := s.db.GetTopic(ctx, topicID)
	if errutil.Is(errutil.NotFound, err) {
		return nil, errutil.E(err).Messagef("topic id='%d' not found", topicID.I64())
	}

	if err != nil {
		return nil, errutil.E(err).Debug("s.db.GetTopic")
	}

	return topic, nil
}
