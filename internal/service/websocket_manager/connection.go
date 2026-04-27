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
	lastPong     time.Time
	ws           *websocket.Conn
	toSend       chan *event_desc.Event
	manager      *WebsocketManager
	cancel       context.CancelFunc
	channel      event_desc.Channel
	userID       domain.UserID
	channelMutex sync.RWMutex
}
