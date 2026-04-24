package text_topic

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) RemoveUserFromTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := s.redis.RemoveUserFromTextTopic(ctx, topicID, userID); err != nil {
		return errutil.E(err).Debug("s.redis.RemoveUserFromTextTopic")
	}

	return nil
}
