package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) UserHasAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	topic, err := s.getTopic(ctx, topicID)
	if err != nil {
		return err
	}

	if _, err = s.db.GetServerParticipant(ctx, topic.ServerID, userID); err != nil {
		return errors.E(err).
			Kind(errors.Permission).
			Message("you don't have an access to the topic")
	}

	return nil
}
