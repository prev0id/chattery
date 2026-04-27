package websocket_manager

import (
	"context"
	"sync"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/set"
)

type session struct {
	ctx         context.Context
	redis       redis
	cancel      context.CancelFunc
	connections set.Set[*Connection]
	userID      domain.UserID
	m           sync.RWMutex
}

func newSession(conn *Connection, userID domain.UserID) *session {
	ctx, cancel := context.WithCancel(context.Background())
	return &session{
		connections: set.NewSet(conn),
		ctx:         ctx,
		cancel:      cancel,
		userID:      userID,
	}
}

func (s *session) addConnection(conn *Connection) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.connections.Empty() {
		go s.startRedisSubscription()
	}
	s.connections.Add(conn)
}

func (s *session) removeConnection(conn *Connection) {
	s.m.Lock()
	defer s.m.Unlock()

	s.connections.Delete(conn)
	if s.connections.Empty() {
		s.cancel()
	}
}

func (s *session) startRedisSubscription() {
	events := make(chan *event_desc.Event)
	go s.redis.SubscribeToUser(s.ctx, s.userID, events)
	go s.handleRedisMessages(events)
}

func (s *session) handleRedisMessages(events chan *event_desc.Event) {
	for event := range events {
		s.broadcastEvent(s.ctx, event)
	}
}

func (s *session) broadcastEvent(ctx context.Context, event *event_desc.Event) {
	s.m.RLock()
	defer s.m.RUnlock()

	for conn := range s.connections {
		conn.SendEvent(ctx, event)
	}
}
