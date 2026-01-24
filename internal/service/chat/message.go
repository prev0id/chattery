package chat

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) PostMessage(ctx context.Context, message *domain.Message) error {
	if _, err := s.db.CreateMessage(ctx, message); err != nil {
		return errors.E(err).Debug("s.db.CreateMessage")
	}
	return nil
}

func (s *Service) ListMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.MessageCursor) ([]*domain.Message, *domain.MessageCursor, error) {
	if cursor == nil {
		return s.firstPageOfMessages(ctx, chatID)
	}
	return s.nextPageOfMessages(ctx, chatID, cursor)
}

func (s *Service) firstPageOfMessages(ctx context.Context, chatID domain.ChatID) ([]*domain.Message, *domain.MessageCursor, error) {
	messages, next, err := s.db.FirstPageOfMessages(ctx, chatID)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.FirstPageOfMessages")
	}
	return messages, next, nil
}

func (s *Service) nextPageOfMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.MessageCursor) ([]*domain.Message, *domain.MessageCursor, error) {
	messages, next, err := s.db.NextPageOfMessages(ctx, chatID, cursor)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.FirstPageOfMessages")
	}
	return messages, next, nil
}
