package server

import (
	"context"
	"slices"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error) {
	servers, err := s.db.GetUserServers(ctx, userID)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.GetUserServers")
	}

	slices.SortFunc(servers, func(lhs, rhs *domain.Server) int {
		return lhs.JoinedAt.Compare(rhs.JoinedAt)
	})

	for _, server := range servers {
		slices.SortFunc(server.Topics, func(lhs, rhs *domain.Topic) int {
			return lhs.CreatedAt.Compare(rhs.CreatedAt)
		})
	}

	return servers, nil
}
