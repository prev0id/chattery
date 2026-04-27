package dm

import (
	"context"

	"chattery/internal/domain"
)

func (s *Service) ValidateAccess(ctx context.Context, userID domain.UserID, dmID domain.DMID) error {
	if _, err := s.db.GetDMParticipant(ctx, dmID, userID); err != nil {
		return err
	}
	return nil
}
