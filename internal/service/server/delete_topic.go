package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) DeleteTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.deleteTopic(ctx, topicID, userID)
	})
}

func (s *Service) deleteTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	if err := s.validateDeleteTopic(ctx, topicID, userID); err != nil {
		return err
	}

	if err := s.db.DeleteTopic(ctx, topicID); err != nil {
		return errors.E(err).Debug("s.db.DeleteTopic")
	}
	return nil
}

func (s *Service) validateDeleteTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	if err := s.validateTopicExists(ctx, topicID); err != nil {
		return err
	}

	topic, err := s.getTopic(ctx, topicID)
	if err != nil {
		return err
	}

	if err := s.validateUserIsOwner(ctx, topic.ServerID, userID); err != nil {
		return err
	}

	return nil
}
