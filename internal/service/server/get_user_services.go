package server

import (
	"context"
	"slices"

	"chattery/internal/domain"
	"chattery/internal/utils/compare"
	"chattery/internal/utils/errors"
)

func (s *Service) GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error) {
	servers, err := s.db.GetUserServers(ctx, userID)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.GetUserServers")
	}

	slices.SortFunc(servers, compare.Servers)

	for _, server := range servers {
		slices.SortFunc(server.Topics, compare.Topics)
	}

	return servers, nil
}
