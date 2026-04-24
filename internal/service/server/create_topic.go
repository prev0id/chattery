package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) CreateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) (domain.TopicID, error) {
	var (
		id  domain.TopicID
		err error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		id, err = s.createTopic(ctx, topic, userID)
		return err
	})

	return id, err
}

func (s *Service) createTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) (domain.TopicID, error) {
	if err := s.validateCreateTopic(ctx, topic, userID); err != nil {
		return 0, err
	}

	id, err := s.db.CreateTopic(ctx, topic)
	if err != nil {
		return 0, errutil.E(err).Debug("s.db.CreateTopic")
	}

	return id, nil
}

func (s *Service) validateCreateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) error {
	if err := s.validateServerExists(ctx, topic.ServerID); err != nil {
		return err
	}

	if err := s.validateUserIsOwner(ctx, topic.ServerID, userID); err != nil {
		return err
	}

	return s.validateTopicNameUnique(ctx, topic)
}
