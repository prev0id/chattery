package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) CreateServer(ctx context.Context, name string, creatorUserID domain.UserID) (domain.ServerID, error) {
	var (
		serverID domain.ServerID
		err      error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		serverID, err = s.db.CreateServer(ctx, name)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateServer")
		}

		participant := &domain.ServerParticipant{
			ServerID: serverID,
			UserID:   creatorUserID,
			Role:     domain.ServerRoleOwner,
		}

		if err = s.db.CreateServerParticipant(ctx, participant); err != nil {
			return errors.E(err).Debug("s.db.CreateServerParticipant")
		}

		return nil
	})

	return serverID, err
}
