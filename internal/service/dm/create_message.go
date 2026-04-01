package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) CreateDMMessage(ctx context.Context, message *domain.DMMessage) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.createDMMessage(ctx, message)
	})
}

func (s *Service) createDMMessage(ctx context.Context, message *domain.DMMessage) error {
	if err := s.validateCreateDMMessage(ctx, message); err != nil {
		return err
	}

	messageID, err := s.db.CreateDMMessage(ctx, message)
	if err != nil {
		return errors.E(err).Debug("s.db.CreateDMMessage")
	}

	if err := s.db.SetLastMessageInDM(ctx, message.DMID, messageID); err != nil {
		return errors.E(err).Debug("s.db.SetLastMessageInDM")
	}

	return nil
}

func (s *Service) validateCreateDMMessage(ctx context.Context, message *domain.DMMessage) error {
	return nil
}
