package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error) {
	dms, err := s.db.UserDMs(ctx, userID)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.UserDMs")
	}

	return dms, nil
}
