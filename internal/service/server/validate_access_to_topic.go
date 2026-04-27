package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) ValidateAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID, targetType domain.TopicType) error {
	topic, err := s.getTopic(ctx, topicID)
	if err != nil {
		return err
	}

	if topic.Type != targetType {
		return errutil.E(err).
			Kind(errutil.InvalidRequest).
			Messagef("topic's type must be %s", targetType.String())
	}

	if _, err = s.db.GetServerParticipant(ctx, topic.ServerID, userID); err != nil {
		return errutil.E(err).
			Kind(errutil.Permission).
			Message("you don't have an access to the topic")
	}

	return nil
}
