package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) FirstPageOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	var (
		messages   []*domain.DMMessage
		nextCursor *domain.DMCursor
		err        error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		messages, nextCursor, err = s.firstPageOfDMMessages(ctx, userID, cursor)
		return err
	})

	return messages, nextCursor, err
}

func (s *Service) firstPageOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	if err := s.validateFirstPageOfDMMessages(ctx, cursor.ChatID, userID); err != nil {
		return nil, nil, err
	}

	cursor.Limit = s.limit

	messages, err := s.db.FirstPageOfDMMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errutil.E(err).Debug("s.db.FirstPageOfDMMessages")
	}

	if len(messages) > 0 {
		lastSeenMessage := messages[0]
		if err := s.db.SetDMLastReadMessage(ctx, cursor.ChatID, userID, lastSeenMessage.ID); err != nil {
			return nil, nil, errutil.E(err).Debug("s.db.SetDMLastReadMessage")
		}
	}

	return messages, s.getNextCursor(cursor.ChatID, messages), nil
}

func (s *Service) validateFirstPageOfDMMessages(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	return s.validateParticipantExists(ctx, dmID, userID)
}
