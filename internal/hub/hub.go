package hub

import (
	"context"
	"sync"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

type dmService interface {
	GetParticipants(ctx context.Context, dmID domain.DMID) ([]domain.UserID, error)
}

type serverService interface {
	AddUserInTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	RemoveUserFromTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	UserHasAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
}

type userCache interface {
	GetByID(userID domain.UserID) (*domain.User, error)
}

type redis interface {
	SubscribeToUser(ctx context.Context, userID domain.UserID, dst chan<- *domain.UserMessage)
	PublishToUser(ctx context.Context, userID domain.UserID, message *domain.UserMessage) error
	RemoveUserFromTextTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
}

type Hub struct {
	redis         redis
	dmService     dmService
	serverService serverService
	user          userCache

	connections     map[domain.UserID]map[*Connection]bool
	userSubscribers map[domain.UserID]chan *domain.UserMessage

	m sync.RWMutex
}

func New(redisAdapter redis, dmSvc dmService, serverSvc serverService, userCache userCache) *Hub {
	return &Hub{
		redis:           redisAdapter,
		dmService:       dmSvc,
		serverService:   serverSvc,
		user:            userCache,
		connections:     make(map[domain.UserID]map[*Connection]bool),
		userSubscribers: make(map[domain.UserID]chan *domain.UserMessage),
	}
}

func (h *Hub) Redis() redis {
	return h.redis
}

func (h *Hub) AddUserInTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := h.serverService.UserHasAccessToTopic(ctx, userID, topicID); err != nil {
		return errors.E(err).Debug("h.serverService.UserHasAccessToTopic")
	}

	if err := h.serverService.AddUserInTextTopic(ctx, userID, topicID); err != nil {
		return errors.E(err).Debug("h.serverService.AddUserInTextTopic")
	}
	return nil
}

func (h *Hub) UpdateUserInTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := h.serverService.AddUserInTextTopic(ctx, userID, topicID); err != nil {
		return errors.E(err).Debug("h.serverService.AddUserInTextTopic")
	}
	return nil
}

func (h *Hub) RemoveUserFromTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	return h.serverService.RemoveUserFromTextTopic(ctx, userID, topicID)
}

func (h *Hub) RegisterConnection(conn *Connection) {
	h.m.Lock()
	defer h.m.Unlock()

	if _, ok := h.connections[conn.userID]; !ok {
		h.connections[conn.userID] = make(map[*Connection]bool)
	}
	h.connections[conn.userID][conn] = true
}

func (h *Hub) UnregisterConnection(conn *Connection) {
	h.m.Lock()
	defer h.m.Unlock()

	if conns, ok := h.connections[conn.userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.connections, conn.userID)
			h.stopUserSubscriber(conn.userID)
		}
	}
}

func (h *Hub) SubscribeUser(ctx context.Context, userID domain.UserID) {
	h.m.Lock()
	defer h.m.Unlock()

	if _, ok := h.userSubscribers[userID]; ok {
		return
	}

	dst := make(chan *domain.UserMessage)
	h.userSubscribers[userID] = dst

	go h.redis.SubscribeToUser(ctx, userID, dst)
	go h.handleUserMessages(userID, dst)
}

func (h *Hub) stopUserSubscriber(userID domain.UserID) {
	if ch, ok := h.userSubscribers[userID]; ok {
		close(ch)
		delete(h.userSubscribers, userID)
	}
}

func (h *Hub) handleUserMessages(userID domain.UserID, src chan *domain.UserMessage) {
	for msg := range src {
		event := h.userMessageToEvent(msg)
		h.sendToUserConnections(userID, event)
	}
}

func (h *Hub) sendToUserConnections(userID domain.UserID, event *Event) {
	h.m.RLock()
	defer h.m.RUnlock()

	conns, ok := h.connections[userID]
	if !ok {
		return
	}

	for conn := range conns {
		select {
		case conn.send <- event:
		default:
		}
	}
}

func (h *Hub) userMessageToEvent(msg *domain.UserMessage) *Event {
	switch msg.Type {
	case domain.UserMessageTypeDM:
		return h.dmToEventMessage(msg.DMMessage)
	case domain.UserMessageTypeTopic:
		return h.topicToEventMessage(msg.TopicMsg)
	default:
		return nil
	}
}

func (h *Hub) dmToEventMessage(msg *domain.DMMessage) *Event {
	sender := UserInfo{}
	if user, err := h.user.GetByID(msg.SenderID); err == nil {
		sender = UserInfo{
			ID:       int64(user.ID),
			Username: user.Username.String(),
			Avatar:   user.AvatarID.String(),
		}
	}

	return &Event{
		Type:        EventMessage,
		ChannelType: ChannelDM,
		ChannelID:   msg.DMID.I64(),
		Message: &EventData{
			ID:        msg.ID.I64(),
			Sender:    sender,
			Text:      msg.Text,
			CreatedAt: render.Timestamp(msg.CreatedAt),
		},
	}
}

func (h *Hub) topicToEventMessage(msg *domain.TopicMessage) *Event {
	sender := UserInfo{}
	if user, err := h.user.GetByID(msg.SenderID); err == nil {
		sender = UserInfo{
			ID:       user.ID.I64(),
			Username: user.Username.String(),
			Avatar:   user.AvatarID.String(),
		}
	}

	return &Event{
		Type:        EventMessage,
		ChannelType: ChannelTextTopic,
		ChannelID:   msg.TopicID.I64(),
		Message: &EventData{
			ID:        msg.ID.I64(),
			Sender:    sender,
			Text:      msg.Text,
			CreatedAt: render.Timestamp(msg.CreatedAt),
		},
	}
}
