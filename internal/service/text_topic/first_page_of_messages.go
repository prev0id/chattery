package text_topic

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) FirstPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error) {
	var (
		messages   []*domain.TopicMessage
		nextCursor *domain.TopicCursor
		err        error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		messages, nextCursor, err = s.firstPageOfMessages(ctx, cursor, userID)
		return err
	})

	return messages, nextCursor, err
}

func (s *Service) firstPageOfMessages(ctx context.Context, cursor *domain.TopicCursor, userID domain.UserID) ([]*domain.TopicMessage, *domain.TopicCursor, error) {
	if err := s.validateFirstPageOfMessages(ctx, cursor.ChatID, userID); err != nil {
		return nil, nil, err
	}

	cursor.Limit = s.limit

	messages, err := s.db.FirstPageOfTopicMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.FirstPageOfTopicMessages")
	}

	return messages, s.getNextCursor(cursor.ChatID, messages), nil
}

func (s *Service) validateFirstPageOfMessages(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	topic, err := s.getTopic(ctx, topicID)
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

func (s *Service) validateParticipantExists(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	_, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetServerParticipant")
	}
	return nil
}

func (s *Service) validateTopicIsText(topic *domain.Topic) error {
	if topic.Type != domain.TopicTypeText {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("the topic id='%d' type must be text", topic.ID.I64())
	}
	return nil
}
