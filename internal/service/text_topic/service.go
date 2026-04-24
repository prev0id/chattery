package text_topic

import (
	"context"
	"time"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

type db interface {
	GetTopic(ctx context.Context, topicID domain.TopicID) (*domain.Topic, error)
	GetServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) (*domain.ServerParticipant, error)
	CreateMessage(ctx context.Context, message *domain.TopicMessage) error
	FirstPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error)
	NextPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error)
}

type redis interface {
	PublishToUser(ctx context.Context, userID domain.UserID, message *domain.UserMessage) error
	AddUserInTextTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
	ListUsersInTextTopic(ctx context.Context, topicID domain.TopicID, threshold time.Duration) ([]domain.UserID, error)
	RemoveUserFromTextTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
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
