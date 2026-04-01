package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) AddParticipant(ctx context.Context, participant *domain.ServerParticipant) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.addParticipant(ctx, participant)
	})
}

func (s *Service) addParticipant(ctx context.Context, participant *domain.ServerParticipant) error {
	if err := s.validateAddParticipant(ctx, participant); err != nil {
		return err
	}

	if err := s.db.CreateServerParticipant(ctx, participant); err != nil {
		return errors.E(err).Debug("s.db.CreateServerParticipant")
	}

	return nil
}

func (s *Service) validateAddParticipant(ctx context.Context, participant *domain.ServerParticipant) error {
	if err := s.validateServerExists(ctx, participant.ServerID); err != nil {
		return err
	}

	if err := s.validateParticipantNotExists(ctx, participant.ServerID, participant.UserID); err != nil {
		return err
	}

	return nil
}
