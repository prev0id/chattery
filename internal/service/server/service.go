package server

import (
	"context"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type db interface {
	CreateServer(ctx context.Context, name string) (domain.ServerID, error)
	UpdateServer(ctx context.Context, serverID domain.ServerID, name string) error
	DeleteServer(ctx context.Context, serverID domain.ServerID) error
	CreateServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID, role domain.ServerRole) error
	DeleteServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error
	CreateTopic(ctx context.Context, serverID domain.ServerID, name string, topicType domain.TopicType) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, topicID domain.TopicID, name string) error
	DeleteTopic(ctx context.Context, topicID domain.TopicID) error
	CreateMessage(ctx context.Context, topicID domain.TopicID, userID domain.UserID, text string) (domain.TopicMessageID, error)
	FirstPageOfTopicMessages(ctx context.Context, cursor domain.TopicCursor) ([]*domain.TopicMessage, error)
	NextPagesOfTopicMessages(ctx context.Context, cursor domain.TopicCursor) ([]*domain.TopicMessage, error)
	GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error)
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

func (s *Service) CreateServer(ctx context.Context, name string, creatorUserID domain.UserID) (domain.ServerID, error) {
	var (
		serverID domain.ServerID
		err      error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		serverID, err = s.db.CreateServer(ctx, name)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateServer")
		}

		if err = s.db.CreateServerParticipant(ctx, serverID, creatorUserID, domain.ServerRoleOwner); err != nil {
			return errors.E(err).Debug("s.db.CreateServerParticipant")
		}

		return nil
	})

	return serverID, err
}

func (s *Service) UpdateServer(ctx context.Context, serverID domain.ServerID, name string) error {
	if err := s.db.UpdateServer(ctx, serverID, name); err != nil {
		return errors.E(err).Debug("s.db.UpdateServer")
	}
	return nil
}

func (s *Service) DeleteServer(ctx context.Context, serverID domain.ServerID) error {
	if err := s.db.DeleteServer(ctx, serverID); err != nil {
		return errors.E(err).Debug("s.db.DeleteServer")
	}
	return nil
}

func (s *Service) AddServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID, role domain.ServerRole) error {
	if err := s.db.CreateServerParticipant(ctx, serverID, userID, role); err != nil {
		return errors.E(err).Debug("s.db.CreateServerParticipant")
	}
	return nil
}

func (s *Service) RemoveServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.db.DeleteServerParticipant(ctx, serverID, userID); err != nil {
		return errors.E(err).Debug("s.db.DeleteServerParticipant")
	}
	return nil
}

func (s *Service) CreateTopic(ctx context.Context, serverID domain.ServerID, name string, topicType domain.TopicType) (*domain.Topic, error) {
	topic, err := s.db.CreateTopic(ctx, serverID, name, topicType)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.CreateTopic")
	}
	return topic, nil
}

func (s *Service) UpdateTopic(ctx context.Context, topicID domain.TopicID, name string) error {
	if err := s.db.UpdateTopic(ctx, topicID, name); err != nil {
		return errors.E(err).Debug("s.db.UpdateTopic")
	}
	return nil
}

func (s *Service) DeleteTopic(ctx context.Context, topicID domain.TopicID) error {
	if err := s.db.DeleteTopic(ctx, topicID); err != nil {
		return errors.E(err).Debug("s.db.DeleteTopic")
	}
	return nil
}

func (s *Service) CreateTopicMessage(ctx context.Context, topicID domain.TopicID, userID domain.UserID, text string) (domain.TopicMessageID, error) {
	messageID, err := s.db.CreateMessage(ctx, topicID, userID, text)
	if err != nil {
		return domain.TopicMessageID(0), errors.E(err).Debug("s.db.CreateMessage")
	}
	return messageID, nil
}

func (s *Service) GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error) {
	servers, err := s.db.GetUserServers(ctx, userID)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.GetUserServers")
	}
	return servers, nil
}

func (s *Service) FirstPageOfTopicMessages(ctx context.Context, cursor domain.TopicCursor) ([]*domain.TopicMessage, error) {
	cursor.Limit = s.limit

	messages, err := s.db.FirstPageOfTopicMessages(ctx, cursor)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.FirstPageOfTopicMessages")
	}
	return messages, nil
}

func (s *Service) NextPagesOfTopicMessages(ctx context.Context, cursor domain.TopicCursor) ([]*domain.TopicMessage, error) {
	cursor.Limit = s.limit

	messages, err := s.db.NextPagesOfTopicMessages(ctx, cursor)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.NextPagesOfTopicMessages")
	}
	return messages, nil
}
