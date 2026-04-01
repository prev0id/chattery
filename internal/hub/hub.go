package hub

import (
	"context"
	"fmt"
	"sync"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type dmService interface {
	UserHasAccessToDM(ctx context.Context, userID domain.UserID, dmID domain.DMID) error
}

type serverService interface {
	UserHasAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
}

type redis interface {
	SubscribeToDM(ctx context.Context, dmID domain.DMID, dst chan<- *domain.DMMessage)
	SubscribeToServerTopic(ctx context.Context, topicID domain.TopicID, dst chan<- *domain.TopicMessage)
}

type Hub struct {
	redis         redis
	dmService     dmService
	serverService serverService

	connections map[domain.UserID]map[*Connection]bool

	rooms map[string]map[domain.UserID]bool

	listeners map[string]context.CancelFunc

	m sync.RWMutex
}

func New(redisAdapter redis, dmSvc dmService, serverSvc serverService) *Hub {
	return &Hub{
		redis:         redisAdapter,
		dmService:     dmSvc,
		serverService: serverSvc,
		connections:   make(map[domain.UserID]map[*Connection]bool),
		rooms:         make(map[string]map[domain.UserID]bool),
		listeners:     make(map[string]context.CancelFunc),
	}
}

func (h *Hub) Join(ctx context.Context, userID domain.UserID, channelType string, channelID int64) error {
	channelKey := fmt.Sprintf("%s:%d", channelType, channelID)

	switch channelType {
	case "dm":
		dmID := domain.DMID(channelID)
		if err := h.dmService.UserHasAccessToDM(ctx, userID, dmID); err != nil {
			return errors.E(err).Kind(errors.Permission).Message("no access to dm")
		}
	case "server":
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
		h.startListener(channelKey, channelType, channelID)
	}

	return nil
}

func (h *Hub) Leave(userID domain.UserID, channelType string, channelID string) {
	channelKey := channelType + ":" + channelID

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

func (h *Hub) GetUsersInChannel(channelKey string) []domain.UserID {
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

func (h *Hub) startListener(channelKey, channelType string, channelID int64) {
	ctx, cancel := context.WithCancel(context.Background())
	h.listeners[channelKey] = cancel

	switch channelType {
	case "dm":
		dmID := domain.DMID(channelID)
		dst := make(chan *domain.DMMessage)
		go h.redis.SubscribeToDM(ctx, dmID, dst)
		go h.handleDMMessages(channelKey, dst)
	case "server":
		topicID := domain.TopicID(channelID)
		dst := make(chan *domain.TopicMessage)
		go h.redis.SubscribeToServerTopic(ctx, topicID, dst)
		go h.handleServerMessages(channelKey, dst)
	}
}

func (h *Hub) stopListener(channelKey string) {
	if cancel, ok := h.listeners[channelKey]; ok {
		cancel()
		delete(h.listeners, channelKey)
	}
}

func (h *Hub) handleDMMessages(channelKey string, src chan *domain.DMMessage) {
	for message := range src {
		h.broadcast(channelKey, message)
	}
}

func (h *Hub) handleServerMessages(channelKey string, src chan *domain.TopicMessage) {
	for message := range src {
		h.broadcast(channelKey, message)
	}
}

func (h *Hub) broadcast(channelKey string, message any) {
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
