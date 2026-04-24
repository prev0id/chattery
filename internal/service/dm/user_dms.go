package dm

import (
	"context"
	"slices"

	"chattery/internal/domain"
	"chattery/internal/utils/compare"
	"chattery/internal/utils/errutil"
)

func (s *Service) UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error) {
	dms, err := s.db.UserDMs(ctx, userID)
	if err != nil {
		return nil, errutil.E(err).Debug("s.db.UserDMs")
	}

	slices.SortFunc(dms, compare.DMs)

	return dms, nil
}
