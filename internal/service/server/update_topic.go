package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) UpdateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.updateTopic(ctx, topic, userID)
	})
}

func (s *Service) updateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) error {
	if err := s.validateUpdateTopic(ctx, topic, userID); err != nil {
		return err
	}

	if err := s.db.UpdateTopic(ctx, topic); err != nil {
		return errors.E(err).Debug("s.db.UpdateTopic")
	}
	return nil
}

func (s *Service) validateUpdateTopic(ctx context.Context, topic *domain.Topic, userID domain.UserID) error {
	if err := s.validateTopicExists(ctx, topic.ID); err != nil {
		return err
	}

	topicFromDB, err := s.getTopic(ctx, topic.ID)
	if err != nil {
		return err
	}

	if err := s.validateUserIsOwner(ctx, topicFromDB.ServerID, userID); err != nil {
		return err
	}

	return nil
}
