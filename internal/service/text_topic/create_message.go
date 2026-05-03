package text_topic

import (
	"context"
	"sync"
	"time"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

func (s *Service) CreateMessage(ctx context.Context, message *domain.TopicMessage) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		return s.createMessage(ctx, message)
	})
}

func (s *Service) createMessage(ctx context.Context, message *domain.TopicMessage) error {
	topic, err := s.getTopic(ctx, message.TopicID)
	if err != nil {
		return err
	}

	if err = s.validateCreateMessage(ctx, message, topic); err != nil {
		return err
	}

	messageID, err := s.db.CreateMessage(ctx, message)
	if err != nil {
		return errutil.E(err).Debug("s.db.CreateMessage")
	}

	message.ID = messageID
	message.CreatedAt = time.Now()

	return s.broadcastMessage(ctx, message, topic.ServerID)
}

func (s *Service) validateCreateMessage(ctx context.Context, message *domain.TopicMessage, topic *domain.Topic) error {
	if err := s.validateParticipantExists(ctx, topic.ServerID, message.SenderID); err != nil {
		return err
	}
	return s.validateTopicIsText(topic)
}

func (s *Service) broadcastMessage(ctx context.Context, message *domain.TopicMessage, serverID domain.ServerID) error {
	channel, payload := s.convertMessageToDesc(message)

	renderedEvent, err := render.Event(event_desc.TypeMessage, channel, payload)
	if err != nil {
		return errutil.E(err).Debug("render.Event")
	}

	participants, err := s.db.GetServerParticipants(ctx, serverID)
	if err != nil {
		return errutil.E(err).Debug("s.db.GetTopicParticipants")
	}

	wg := sync.WaitGroup{}
	for _, participant := range participants {
		wg.Go(func() {
			s.sendEvent(ctx, participant.UserID, renderedEvent)
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

func (s *Service) convertMessageToDesc(message *domain.TopicMessage) (event_desc.Channel, *event_desc.MessagePayload) {
	payload := &event_desc.MessagePayload{
		ID:        message.ID.I64(),
		Text:      message.Text,
		CreatedAt: render.Timestamp(message.CreatedAt),
		Sender:    s.getUserInfo(message.SenderID),
	}
	channel := event_desc.Channel{
		Type: event_desc.ChannelTextTopic,
		ID:   message.TopicID.I64(),
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
		Avatar:   render.AvatarURL(user.Username),
		ID:       user.ID.I64(),
	}
}
