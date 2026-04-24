package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) UpdateServer(ctx context.Context, serverID domain.ServerID, name string, userID domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.updateServer(ctx, serverID, name, userID)
	})
}

func (s *Service) updateServer(ctx context.Context, serverID domain.ServerID, name string, userID domain.UserID) error {
	if err := s.validateServerUpdate(ctx, serverID, userID); err != nil {
		return err
	}

	updatedServer := &domain.Server{
		ID:   serverID,
		Name: name,
	}

	if err := s.db.UpdateServer(ctx, updatedServer); err != nil {
		return errutil.E(err).Debug("s.db.UpdateServer")
	}

	return nil
}

func (s *Service) validateServerUpdate(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.validateServerExists(ctx, serverID); err != nil {
		return err
	}

	return s.validateUserIsOwner(ctx, serverID, userID)
}
