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

func newSession(redisAdapter redis, conn *Connection, userID domain.UserID) *session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &session{
		connections: set.NewSet(conn),
		ctx:         ctx,
		redis:       redisAdapter,
		cancel:      cancel,
		userID:      userID,
	}
	go s.startRedisSubscription()
	return s
}

func (s *session) addConnection(conn *Connection) {
	s.m.Lock()
	defer s.m.Unlock()

	s.connections.Add(conn)
}

func (s *session) removeConnection(conn *Connection) bool {
	s.m.Lock()
	defer s.m.Unlock()

	s.connections.Delete(conn)
	if s.connections.Empty() {
		s.cancel()
		return true
	}
	return false
}

func (s *session) startRedisSubscription() {
	events := make(chan *event_desc.Event)
	go s.redis.SubscribeToUser(s.ctx, s.userID, events)
	s.handleRedisMessages(events)
}

func (s *session) handleRedisMessages(events chan *event_desc.Event) {
	for event := range events {
		s.broadcastEvent(s.ctx, event)
	}
}

func (s *session) broadcastEvent(ctx context.Context, event *event_desc.Event) {
	s.m.RLock()
	connections := make([]*Connection, 0, len(s.connections))
	for conn := range s.connections {
		connections = append(connections, conn)
	}
	s.m.RUnlock()

	for _, conn := range connections {
		if !conn.shouldReceive(event) {
			continue
		}
		conn.SendEvent(ctx, event)
	}
}
