package dm

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
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
		return errutil.E(err).Debug("s.db.CreateDMMessage")
	}

	if err = s.db.SetLastMessageInDM(ctx, message.DMID, messageID); err != nil {
		return errutil.E(err).Debug("s.db.SetLastMessageInDM")
	}

	message.ID = messageID

	participants, err := s.db.GetDMParticipants(ctx, message.DMID)
	if err != nil {
		return errutil.E(err).Debug("s.db.GetDMParticipants")
	}

	userMsg := &domain.UserMessage{
		Type:      domain.UserMessageTypeDM,
		ChannelID: message.DMID.I64(),
		DMMessage: message,
	}

	for _, participantID := range participants {
		if err := s.redis.PublishToUser(ctx, participantID, userMsg); err != nil {
			return errutil.E(err).Debug("s.redis.PublishToUser")
		}
	}

	return nil
}

func (s *Service) validateCreateDMMessage(ctx context.Context, message *domain.DMMessage) error {
	return s.validateParticipantExists(ctx, message.DMID, message.SenderID)
}
