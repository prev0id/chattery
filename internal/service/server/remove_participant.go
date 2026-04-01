package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) RemoveParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.removeParticipant(ctx, serverID, userID)
	})
}

func (s *Service) removeParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.validateRemoveParticipant(ctx, serverID, userID); err != nil {
		return err
	}

	if err := s.db.DeleteServerParticipant(ctx, serverID, userID); err != nil {
		return errors.E(err).Debug("s.db.DeleteServerParticipant")
	}

	return nil
}

func (s *Service) validateRemoveParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	if err := s.validateServerExists(ctx, serverID); err != nil {
		return err
	}

	if err := s.validateParticipantExists(ctx, serverID, userID); err != nil {
		return err
	}

	if err := s.validateUserIsNotOwner(ctx, serverID, userID); err != nil {
		return errors.E(err).Message("owner can't leave the server")
	}

	return nil
}
