package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) DeleteServer(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.deleteServer(ctx, serverID, userID)
	})
}

func (s *Service) deleteServer(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.validateDeleteServer(ctx, serverID, userID); err != nil {
		return err
	}

	if err := s.db.DeleteServer(ctx, serverID); err != nil {
		return errutil.E(err).Debug("s.db.DeleteServer")
	}

	return nil
}

func (s *Service) validateDeleteServer(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.validateServerExists(ctx, serverID); err != nil {
		return err
	}

	return s.validateUserIsOwner(ctx, serverID, userID)
}
