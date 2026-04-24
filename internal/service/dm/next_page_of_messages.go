package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) NextPagesOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	var (
		messages   []*domain.DMMessage
		nextCursor *domain.DMCursor
		err        error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		messages, nextCursor, err = s.nextPagesOfDMMessages(ctx, userID, cursor)
		return err
	})

	return messages, nextCursor, err
}

func (s *Service) nextPagesOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	if err := s.validateNextPagesOfDMMessages(ctx, cursor.ChatID, userID); err != nil {
		return nil, nil, err
	}

	cursor.Limit = s.limit

	messages, err := s.db.NextPagesOfDMMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errutil.E(err).Debug("s.db.NextPagesOfDMMessages")
	}

	return messages, s.getNextCursor(cursor.ChatID, messages), nil
}

func (s *Service) validateNextPagesOfDMMessages(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	return s.validateParticipantExists(ctx, dmID, userID)
}
