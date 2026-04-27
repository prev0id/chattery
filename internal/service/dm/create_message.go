package dm

import (
	"context"
	"sync"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
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

	return s.broadcastMessage(ctx, message)
}

func (s *Service) validateCreateDMMessage(ctx context.Context, message *domain.DMMessage) error {
	return s.validateParticipantExists(ctx, message.DMID, message.SenderID)
}

func (s *Service) broadcastMessage(ctx context.Context, message *domain.DMMessage) error {
	channel, payload := s.convertMessageToDesc(message)

	renderedEvent, err := render.Event(event_desc.TypeMessage, channel, payload)
	if err != nil {
		return errutil.E(err).Debug("render.Event")
	}

	participants, err := s.db.GetDMParticipants(ctx, message.DMID)
	if err != nil {
		return errutil.E(err).Debug("s.db.GetDMParticipants")
	}

	wg := sync.WaitGroup{}
	for _, participant := range participants {
		wg.Go(func() {
			s.sendEvent(ctx, participant, renderedEvent)
		})
	}
	wg.Wait()

	return nil
}

func (s *Service) sendEvent(ctx context.Context, to domain.UserID, event []byte) {
	if err := s.redis.PublishToUser(ctx, to, event); err != nil {
		logger.ErrorCtx(ctx, err, "s.redis.PublishToUser")
	}
}

func (s *Service) convertMessageToDesc(message *domain.DMMessage) (event_desc.Channel, *event_desc.MessagePayload) {
	payload := &event_desc.MessagePayload{
		ID:        message.ID.I64(),
		Text:      message.Text,
		CreatedAt: render.Timestamp(message.CreatedAt),
		Sender:    s.getUserInfo(message.SenderID),
	}
	channel := event_desc.Channel{
		Type: event_desc.ChannelTextTopic,
		ID:   message.DMID.I64(),
	}
	return channel, payload
}

func (s *Service) getUserInfo(userID domain.UserID) event_desc.UserInfo {
	user, err := s.user.GetByID(userID)
	if err != nil {
		logger.Error(err, "s.user.GetByID")
		return event_desc.UserInfo{}
	}

	return event_desc.UserInfo{
		Username: user.Username.String(),
		Avatar:   "/v1/image/" + user.Username.String() + ".png",
		ID:       user.ID.I64(),
	}
}
