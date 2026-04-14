package server

import (
	"context"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

const userActivityTTL = 90 * time.Second

func (s *Service) CreateMessage(ctx context.Context, message *domain.TopicMessage) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.createMessage(ctx, message)
	})
}

func (s *Service) createMessage(ctx context.Context, message *domain.TopicMessage) error {
	if err := s.validateCreateMessage(ctx, message); err != nil {
		return err
	}

	if err := s.db.CreateMessage(ctx, message); err != nil {
		return errors.E(err).Debug("s.db.CreateMessage")
	}

	participants, err := s.redis.ListUsersInTextTopic(ctx, message.TopicID, userActivityTTL)
	if err != nil {
		return errors.E(err).Debug("s.db.GetTopicParticipants")
	}

	userMsg := &domain.UserMessage{
		Type:      domain.UserMessageTypeTopic,
		ChannelID: message.TopicID.I64(),
		TopicMsg:  message,
	}

	for _, participantID := range participants {
		if err := s.redis.PublishToUser(ctx, participantID, userMsg); err != nil {
			errors.E(err).Debug("s.redis.PublishToUser")
		}
	}

	return nil
}

func (s *Service) validateCreateMessage(ctx context.Context, message *domain.TopicMessage) error {
	topic, err := s.GetTopic(ctx, message.TopicID)
	if err != nil {
		return err
	}

	if err := s.validateParticipantExists(ctx, topic.ServerID, message.SenderID); err != nil {
		return err
	}

	if err := s.validateTopicIsText(topic); err != nil {
		return err
	}

	return nil
}
