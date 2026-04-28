package server

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

func (s *Service) SearchServersNotJoined(ctx context.Context, userID domain.UserID, query string) ([]*domain.Server, error) {
	userServers, err := s.db.GetUserServers(ctx, userID)
	if err != nil {
		return nil, errutil.E(err).Debug("s.db.GetUserServers")
	}

	joined := set.NewSet(sliceutil.Map(userServers, func(server *domain.Server) domain.ServerID {
		return server.ID
	})...)

	query = strings.ToLower(query)
	filtered := sliceutil.Filter(s.cache.List(), func(server *domain.Server) bool {
		return server != nil &&
			!joined.Contains(server.ID) &&
			(query == "" || strings.Contains(strings.ToLower(server.Name), query))
	})

	slices.SortFunc(filtered, compare.ServersByName)

	return filtered, nil
}
