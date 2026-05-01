package websocket_manager

import (
	"context"
	"sync"
	"time"

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

type voiceService interface {
	Join(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	Leave(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	HandleSignal(ctx context.Context, userID domain.UserID, event *event_desc.Event) error
}

type WebsocketManager struct {
	redis  redis
	dm     dmService
	server serverService
	voice  voiceService

	sessions      map[domain.UserID]*session
	sessionsMutex sync.RWMutex
}

func New(redisAdapter redis, dm dmService, server serverService, voice voiceService) *WebsocketManager {
	return &WebsocketManager{
		redis:    redisAdapter,
		dm:       dm,
		server:   server,
		voice:    voice,
		sessions: make(map[domain.UserID]*session),
	}
}

func (m *WebsocketManager) NewConnection(ctx context.Context, userID domain.UserID, ws *websocket.Conn) *Connection {
	connCtx, cancel := context.WithCancel(ctx)
	conn := &Connection{
		ctx:      connCtx,
		lastPong: time.Now(),
		manager:  m,
		userID:   userID,
		ws:       ws,
		toSend:   make(chan *event_desc.Event, 256),
		cancel:   cancel,
	}

	m.registerConnection(conn)
	return conn
}

func (m *WebsocketManager) registerConnection(conn *Connection) {
	m.sessionsMutex.Lock()
	defer m.sessionsMutex.Unlock()

	if _, ok := m.sessions[conn.userID]; !ok {
		m.sessions[conn.userID] = newSession(m.redis, conn, conn.userID)
		return
	}
	m.sessions[conn.userID].addConnection(conn)
}

func (m *WebsocketManager) unregisterConnection(conn *Connection) {
	m.sessionsMutex.Lock()
	defer m.sessionsMutex.Unlock()

	if session, ok := m.sessions[conn.userID]; ok {
		if session.removeConnection(conn) {
			delete(m.sessions, conn.userID)
		}
	}
}
