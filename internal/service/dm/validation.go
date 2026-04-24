package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) validateParticipantExists(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	_, err := s.db.GetDMParticipant(ctx, dmID, userID)
	if errutil.Is(errutil.NotFound, err) {
		return errutil.E(err).
			Kind(errutil.NotFound).
			Messagef("user id='%d' not a participant of dm id='%d'", userID.I64(), dmID.I64())
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetDMParticipant")
	}
	return nil
}

func (s *Service) validateDMNotExistsBetweenUsers(ctx context.Context, userID1, userID2 domain.UserID) error {
	_, err := s.db.GetDMBetweenUsers(ctx, userID1, userID2)
	if errutil.Is(errutil.NotFound, err) {
		return nil
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetDMBetweenUsers")
	}

	return errutil.E().
		Kind(errutil.Exist).
		Messagef("dm already exists between users id='%d' and id='%d'", userID1.I64(), userID2.I64())
}

func (*Service) validateDifferentUsers(userID1, userID2 domain.UserID) error {
	if userID1 == userID2 {
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Message("cannot create dm with yourself")
	}
	return nil
}
