package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coder/websocket"

	"chattery/internal/domain"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

type EventType string

const (
	EventJoin    EventType = "join"
	EventLeave   EventType = "leave"
	EventMessage EventType = "message"
	EventError   EventType = "error"
)

type ChannelType string

const (
	ChannelDM     ChannelType = "dm"
	ChannelServer ChannelType = "server"
)

type wsEvent struct {
	Type        EventType   `json:"type"`
	ChannelType ChannelType `json:"channel_type,omitempty"`
	ChannelID   int64       `json:"channel_id,omitempty"`
	Message     any         `json:"message,omitempty"`
	Error       string      `json:"error,omitempty"`
}

type Connection struct {
	hub        *Hub
	userID     domain.UserID
	ws         *websocket.Conn
	send       chan any
	active     bool
	channelKey string
}

func NewConnection(hub *Hub, userID domain.UserID, ws *websocket.Conn) *Connection {
	return &Connection{
		hub:    hub,
		userID: userID,
		ws:     ws,
		send:   make(chan any, 256),
		active: false,
	}
}

func (c *Connection) GetUserID() domain.UserID {
	return c.userID
}

func (c *Connection) GetChannelKey() string {
	return c.channelKey
}

func (c *Connection) ReadPump(ctx context.Context) {
	defer func() {
		c.hub.Leave(c.userID, c.getChannelType(), c.getChannelID())
		c.hub.UnregisterConnection(c)
		c.ws.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, rawEvent, err := c.ws.Read(ctx)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Error(err, "[connection] ws.Read")
			}
			return
		}

		var event wsEvent
		if err := json.Unmarshal(rawEvent, &event); err != nil {
			c.sendError("invalid json format")
			continue
		}

		switch event.Type {
		case EventJoin:
			if event.ChannelType == "" || event.ChannelID == 0 {
				c.sendError("channel_type and channel_id required")
				continue
			}
			if c.active {
				c.sendError("already joined a channel")
				continue
			}
			if err := c.hub.Join(ctx, c.userID, string(event.ChannelType), event.ChannelID); err != nil {
				c.sendError(err.Error())
				continue
			}
			c.active = true
			c.channelKey = fmt.Sprintf("%s:%d", event.ChannelType, event.ChannelID)

		case EventLeave:
			if !c.active {
				continue
			}
			c.hub.Leave(c.userID, c.getChannelType(), c.getChannelID())
			c.active = false
			c.channelKey = ""

		default:
			c.sendError("unknown event type")
		}
	}
}

func (c *Connection) WritePump() {
	defer c.ws.Close(websocket.StatusNormalClosure, "")

	for {
		message, ok := <-c.send
		if !ok {
			c.ws.Write(context.Background(), websocket.MessageText, []byte(`{"type":"error","message":"connection closed"}`))
			return
		}

		bytes, err := render.JsonBytes(message)
		if err != nil {
			logger.Error(err, "[connection] render.JsonBytes")
			continue
		}

		if err := c.ws.Write(context.Background(), websocket.MessageText, bytes); err != nil {
			logger.Error(err, "[connection] ws.Write")
			return
		}
	}
}

func (c *Connection) sendError(msg string) {
	event := wsEvent{
		Type:  EventError,
		Error: msg,
	}
	bytes, err := render.JsonBytes(event)
	if err != nil {
		logger.Error(err, "[connection] render.JsonBytes error")
		return
	}
	c.ws.Write(context.Background(), websocket.MessageText, bytes)
}

func (c *Connection) getChannelType() string {
	if c.channelKey == "" {
		return ""
	}
	for i, r := range c.channelKey {
		if r == ':' {
			return c.channelKey[:i]
		}
	}
	return ""
}

func (c *Connection) getChannelID() string {
	if c.channelKey == "" {
		return ""
	}
	for i, r := range c.channelKey {
		if r == ':' {
			return c.channelKey[i+1:]
		}
	}
	return ""
}

func (c *Connection) SetActive(channelKey string) {
	c.active = true
	c.channelKey = channelKey
}

func (c *Connection) IsActive() bool {
	return c.active
}
