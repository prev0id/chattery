package websocket_manager

import (
	"context"
	"sync"

	"github.com/coder/websocket"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
)

type redis interface {
	SubscribeToUser(ctx context.Context, userID domain.UserID, dst chan<- *event_desc.Event)
	PublishToUser(ctx context.Context, userID domain.UserID, event []byte) error
}

type dmService interface {
	ValidateAccess(ctx context.Context, userID domain.UserID, dmID domain.DMID) error
}

type serverService interface {
	ValidateAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID, topicType domain.TopicType) error
}

type WebsocketManager struct {
	redis  redis
	dm     dmService
	server serverService

	sessions      map[domain.UserID]*session
	sessionsMutex sync.RWMutex
}

func New(redisAdapter redis, dm dmService, server serverService) *WebsocketManager {
	return &WebsocketManager{
		redis:    redisAdapter,
		dm:       dm,
		server:   server,
		sessions: make(map[domain.UserID]*session),
	}
}

func (m *WebsocketManager) NewConnection(userID domain.UserID, ws *websocket.Conn) *Connection {
	conn := &Connection{
		manager: m,
		userID:  userID,
		ws:      ws,
		toSend:  make(chan *event_desc.Event, 256),
		cancel:  func() {},
	}

	m.registerConnection(conn)
	return conn
}

func (m *WebsocketManager) registerConnection(conn *Connection) {
	m.sessionsMutex.Lock()
	defer m.sessionsMutex.Unlock()

	if _, ok := m.sessions[conn.userID]; !ok {
		m.sessions[conn.userID] = newSession(conn, conn.userID)
		return
	}
	m.sessions[conn.userID].addConnection(conn)
}

func (m *WebsocketManager) unregisterConnection(conn *Connection) {
	m.sessionsMutex.Lock()
	defer m.sessionsMutex.Unlock()

	if session, ok := m.sessions[conn.userID]; ok {
		session.removeConnection(conn)
		if session.connections.Empty() {
			delete(m.sessions, conn.userID)
		}
	}
}
