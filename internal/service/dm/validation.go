package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) validateParticipantExists(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	_, err := s.db.GetDMParticipant(ctx, dmID, userID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).
			Kind(errors.NotFound).
			Messagef("user id='%d' not a participant of dm id='%d'", userID.I64(), dmID.I64())
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetDMParticipant")
	}
	return nil
}

func (s *Service) validateParticipantNotExists(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	_, err := s.db.GetDMParticipant(ctx, dmID, userID)
	if errors.Is(errors.NotFound, err) {
		return nil
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetDMParticipant")
	}

	return errors.E().
		Kind(errors.Exist).
		Messagef("user id='%d' already a participant of dm id='%d'", userID.I64(), dmID.I64())
}

func (s *Service) validateDMNotExistsBetweenUsers(ctx context.Context, userID1, userID2 domain.UserID) error {
	_, err := s.db.GetDMBetweenUsers(ctx, userID1, userID2)
	if errors.Is(errors.NotFound, err) {
		return nil
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetDMBetweenUsers")
	}

	return errors.E().
		Kind(errors.Exist).
		Messagef("dm already exists between users id='%d' and id='%d'", userID1.I64(), userID2.I64())
}

func (s *Service) validateDifferentUsers(userID1, userID2 domain.UserID) error {
	if userID1 == userID2 {
		return errors.E().
			Kind(errors.InvalidRequest).
			Message("cannot create dm with yourself")
	}
	return nil
}
