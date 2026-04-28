package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) MarkDMMessageRead(ctx context.Context, userID domain.UserID, dmID domain.DMID, messageID domain.DMMessageID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.markDMMessageRead(ctx, userID, dmID, messageID)
	})
}

func (s *Service) markDMMessageRead(ctx context.Context, userID domain.UserID, dmID domain.DMID, messageID domain.DMMessageID) error {
	if err := s.validateMarkDMMessageRead(ctx, userID, dmID, messageID); err != nil {
		return err
	}

	if err := s.db.SetDMLastReadMessage(ctx, dmID, userID, messageID); err != nil {
		return errutil.E(err).Debug("s.db.SetDMLastReadMessage")
	}

	return nil
}
func (s *Service) validateMarkDMMessageRead(ctx context.Context, userID domain.UserID, dmID domain.DMID, messageID domain.DMMessageID) error {
	if err := s.validateParticipantExists(ctx, dmID, userID); err != nil {
		return err
	}

	return s.validateDMMessageExists(ctx, dmID, messageID)
}
