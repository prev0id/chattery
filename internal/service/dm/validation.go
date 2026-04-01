package dm

import (
	"context"

	"chattery/internal/domain"
)

func (s *Service) validateDMExists(ctx context.Context, dmID domain.DMID) error {
	return nil
}

func (s *Service) validateParticipantExists(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	return nil
}

func (s *Service) validateParticipantNotExists(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	return nil
}
