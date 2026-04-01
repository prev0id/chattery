package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) CreateDM(ctx context.Context, participant1, participant2 domain.UserID) (domain.DMID, error) {
	var (
		dmID domain.DMID
		err  error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		dmID, err = s.createDM(ctx, participant1, participant2)
		return err
	})

	return dmID, err
}

func (s *Service) createDM(ctx context.Context, participant1, participant2 domain.UserID) (domain.DMID, error) {
	if err := s.validateCreateDM(ctx, participant1, participant2); err != nil {
		return 0, err
	}

	dmID, err := s.db.CreateDM(ctx)
	if err != nil {
		return 0, errors.E(err).Debug("s.db.CreateDM")
	}

	if err := s.db.CreateDMParticipant(ctx, dmID, participant1); err != nil {
		return 0, errors.E(err).Debug("s.db.CreateDMParticipant")
	}

	if err := s.db.CreateDMParticipant(ctx, dmID, participant2); err != nil {
		return 0, errors.E(err).Debug("s.db.CreateDMParticipant")
	}

	return dmID, nil
}

func (s *Service) validateCreateDM(ctx context.Context, participant1, participant2 domain.UserID) error {
	return nil
}
