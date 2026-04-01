package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) GetTopic(ctx context.Context, topicID domain.TopicID) (*domain.Topic, error) {
	topic, err := s.db.GetTopic(ctx, topicID)
	if errors.Is(errors.NotFound, err) {
		return nil, errors.E(err).Messagef("topic id='%d' not found", topicID.I64())
	}

	if err != nil {
		return nil, errors.E(err).Debug("s.db.GetTopic")
	}

	return topic, nil
}
