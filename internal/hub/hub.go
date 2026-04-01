package hub

import (
	"context"
	"sync"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type dmService interface {
	UserHasAccessToDM(ctx context.Context, userID domain.UserID, dmID domain.DMID) error
}

type serverService interface {
	UserHasAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
}

type userService interface {
	GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
}

type redis interface {
	SubscribeToDM(ctx context.Context, dmID domain.DMID, dst chan<- *domain.DMMessage)
	SubscribeToServerTopic(ctx context.Context, topicID domain.TopicID, dst chan<- *domain.TopicMessage)
}

type Hub struct {
	redis         redis
	dmService     dmService
	serverService serverService
	user          userService

	connections map[domain.UserID]map[*Connection]bool

	rooms map[ChannelKey]map[domain.UserID]bool

	listeners map[ChannelKey]context.CancelFunc

	m sync.RWMutex
}

func New(redisAdapter redis, dmSvc dmService, serverSvc serverService, userSvc userService) *Hub {
	return &Hub{
		redis:         redisAdapter,
		dmService:     dmSvc,
		serverService: serverSvc,
		user:          userSvc,
		connections:   make(map[domain.UserID]map[*Connection]bool),
		rooms:         make(map[ChannelKey]map[domain.UserID]bool),
		listeners:     make(map[ChannelKey]context.CancelFunc),
	}
}

func (h *Hub) Join(ctx context.Context, userID domain.UserID, channelType ChannelType, channelID int64) error {
	channelKey := ChannelKey{Type: channelType, ID: channelID}

	switch channelType {
	case ChannelDM:
		dmID := domain.DMID(channelID)
		if err := h.dmService.UserHasAccessToDM(ctx, userID, dmID); err != nil {
			return errors.E(err).Kind(errors.Permission).Message("no access to dm")
		}
	case ChannelServer:
		topicID := domain.TopicID(channelID)
		if err := h.serverService.UserHasAccessToTopic(ctx, userID, topicID); err != nil {
			return errors.E(err).Kind(errors.Permission).Message("no access to topic")
		}
	default:
		return errors.E().Kind(errors.InvalidRequest).Message("invalid channel type")
	}

	h.m.Lock()
	defer h.m.Unlock()

	if _, ok := h.rooms[channelKey]; !ok {
		h.rooms[channelKey] = make(map[domain.UserID]bool)
	}
	h.rooms[channelKey][userID] = true

	if len(h.rooms[channelKey]) == 1 {
		h.startListener(channelKey)
	}

	return nil
}

func (h *Hub) Leave(userID domain.UserID, channelType ChannelType, channelID int64) {
	channelKey := ChannelKey{Type: channelType, ID: channelID}

	h.m.Lock()
	defer h.m.Unlock()

	if users, ok := h.rooms[channelKey]; ok {
		delete(users, userID)
		if len(users) == 0 {
			delete(h.rooms, channelKey)
			h.stopListener(channelKey)
		}
	}
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
		}
	}
}

func (h *Hub) GetUsersInChannel(channelKey ChannelKey) []domain.UserID {
	h.m.RLock()
	defer h.m.RUnlock()

	users, ok := h.rooms[channelKey]
	if !ok {
		return nil
	}

	result := make([]domain.UserID, 0, len(users))
	for userID := range users {
		result = append(result, userID)
	}
	return result
}

func (h *Hub) GetUserConnections(userID domain.UserID) []*Connection {
	h.m.RLock()
	defer h.m.RUnlock()

	conns, ok := h.connections[userID]
	if !ok {
		return nil
	}

	result := make([]*Connection, 0, len(conns))
	for conn := range conns {
		result = append(result, conn)
	}
	return result
}

func (h *Hub) startListener(channelKey ChannelKey) {
	ctx, cancel := context.WithCancel(context.Background())
	h.listeners[channelKey] = cancel

	switch channelKey.Type {
	case ChannelDM:
		dst := make(chan *domain.DMMessage)
		go h.redis.SubscribeToDM(ctx, channelKey.DMID(), dst)
		go h.handleDMMessages(ctx, channelKey, dst)
	case ChannelServer:
		dst := make(chan *domain.TopicMessage)
		go h.redis.SubscribeToServerTopic(ctx, channelKey.TopicID(), dst)
		go h.handleServerMessages(ctx, channelKey, dst)
	}
}

func (h *Hub) stopListener(channelKey ChannelKey) {
	if cancel, ok := h.listeners[channelKey]; ok {
		cancel()
		delete(h.listeners, channelKey)
	}
}

func (h *Hub) handleDMMessages(ctx context.Context, channelKey ChannelKey, src chan *domain.DMMessage) {
	for message := range src {
		event := h.dmToEventMessage(ctx, message)
		h.broadcast(channelKey, event)
	}
}

func (h *Hub) handleServerMessages(ctx context.Context, channelKey ChannelKey, src chan *domain.TopicMessage) {
	for message := range src {
		event := h.topicToEventMessage(ctx, message)
		h.broadcast(channelKey, event)
	}
}

func (h *Hub) broadcast(channelKey ChannelKey, message any) {
	users := h.GetUsersInChannel(channelKey)

	h.m.RLock()
	defer h.m.RUnlock()

	for _, userID := range users {
		conns, ok := h.connections[userID]
		if !ok {
			continue
		}
		for conn := range conns {
			select {
			case conn.send <- message:
			default:
			}
		}
	}
}

func (h *Hub) dmToEventMessage(ctx context.Context, msg *domain.DMMessage) *Event {
	sender := UserInfo{}
	if user, err := h.user.GetByID(ctx, msg.SenderID); err == nil {
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
			CreatedAt: formatTimestamp(msg.CreatedAt),
		},
	}
}

func (h *Hub) topicToEventMessage(ctx context.Context, msg *domain.TopicMessage) *Event {
	sender := UserInfo{}
	if user, err := h.user.GetByID(ctx, msg.SenderID); err == nil {
		sender = UserInfo{
			ID:       user.ID.I64(),
			Username: user.Username.String(),
			Avatar:   user.AvatarID.String(),
		}
	}

	return &Event{
		Type:        EventMessage,
		ChannelType: ChannelServer,
		ChannelID:   msg.TopicID.I64(),
		Message: &EventData{
			ID:        msg.ID.I64(),
			Sender:    sender,
			Text:      msg.Text,
			CreatedAt: formatTimestamp(msg.CreatedAt),
		},
	}
}

func formatTimestamp(t time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	msgDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	if msgDate.Equal(today) {
		return "Today, " + t.Format("15:04")
	}
	if msgDate.Equal(yesterday) {
		return "Yesterday, " + t.Format("15:04")
	}
	return t.Format("Jan 2, 15:04")
}
