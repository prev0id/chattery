package text_topic

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) AddUserInTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := s.redis.AddUserInTextTopic(ctx, topicID, userID); err != nil {
		return errors.E(err).Debug("s.redis.AddUserInTextTopic")
	}

	return nil
}
