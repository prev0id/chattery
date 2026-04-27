package websocket_manager

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
)

const (
	pingInterval = 15 * time.Second
	pongTimeout  = 60 * time.Second
)

type Connection struct {
	ctx          context.Context
	lastPong     time.Time
	ws           *websocket.Conn
	toSend       chan *event_desc.Event
	manager      *WebsocketManager
	cancel       context.CancelFunc
	channel      event_desc.Channel
	userID       domain.UserID
	channelMutex sync.RWMutex
}

func (c *Connection) shouldReceive(event *event_desc.Event) bool {
	if event == nil || event.Type != event_desc.TypeMessage {
		return true
	}

	c.channelMutex.RLock()
	defer c.channelMutex.RUnlock()

	return c.channel == event.Channel
}
