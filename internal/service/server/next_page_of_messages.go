package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) NextPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error) {
	var (
		messages   []*domain.TopicMessage
		nextCursor *domain.TopicCursor
		err        error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		messages, nextCursor, err = s.nextPageOfMessages(ctx, cursor, userID)
		return err
	})

	return messages, nextCursor, err
}

func (s *Service) nextPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error) {
	if err := s.validateNextPageOfMessages(ctx, cursor.ChatID, userID); err != nil {
		return nil, nil, err
	}

	cursor.Limit = s.limit

	messages, err := s.db.NextPageOfTopicMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.FirstPageOfTopicMessages")
	}

	return messages, s.getNextCursor(cursor.ChatID, messages), nil
}

func (s *Service) validateNextPageOfMessages(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	topic, err := s.GetTopic(ctx, topicID)
	if err != nil {
		return err
	}

	if err := s.validateParticipantExists(ctx, topic.ServerID, userID); err != nil {
		return err
	}

	if err := s.validateTopicIsText(topic); err != nil {
		return err
	}

	return nil
}
