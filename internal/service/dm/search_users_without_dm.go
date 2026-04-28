package dm

import (
	"context"
	"slices"
	"strings"

	"chattery/internal/domain"
	"chattery/internal/utils/compare"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/set"
	"chattery/internal/utils/sliceutil"
)

func (s *Service) SearchUsersWithoutDM(ctx context.Context, userID domain.UserID, query string) ([]*domain.User, error) {
	knownUserIDs, err := s.db.GetUserDMParticipantIDs(ctx, userID)
	if err != nil {
		return nil, errutil.E(err).Debug("s.db.GetUserDMParticipantIDs")
	}

	knownUsers := set.NewSet(knownUserIDs...)
	knownUsers.Add(userID)

	query = strings.TrimSpace(strings.ToLower(query))

	users := filterKnownUsers(s.user.List(), knownUsers, query)

	slices.SortFunc(users, compare.UsersByUsername)

	return users, nil
}

func filterKnownUsers(users []*domain.User, knownUsers set.Set[domain.UserID], query string) []*domain.User {
	return sliceutil.Filter(users, func(user *domain.User) bool {
		return !knownUsers.Contains(user.ID) && strings.Contains(strings.ToLower(user.Username.String()), query)
	})
}
