package websocket_manager

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/coder/websocket"

	"chattery/internal/domain"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

const pingInterval = 15 * time.Second

type Connection struct {
	manager     *WebsocketManager
	userID      domain.UserID
	channelType ChannelType
	channelID   int64
	ws          *websocket.Conn
	send        chan *Event
	active      bool
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewConnection(manager *WebsocketManager, userID domain.UserID, channelType ChannelType, channelID int64, ws *websocket.Conn, ctx context.Context) *Connection {
	ctx, cancel := context.WithCancel(ctx)
	return &Connection{
		manager:     manager,
		userID:      userID,
		channelType: channelType,
		channelID:   channelID,
		ws:          ws,
		send:        make(chan *Event, 256),
		active:      true,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (c *Connection) GetUserID() domain.UserID {
	return c.userID
}

func (c *Connection) ReadPump(ctx context.Context) {
	defer func() {
		c.cancel()
		c.manager.UnregisterConnection(c)
		if c.channelType == ChannelTextTopic {
			c.manager.LeaveFromTextTopic(c.ctx, c.userID, domain.TopicID(c.channelID))
		}
		close(c.send)
		c.ws.Close(websocket.StatusNormalClosure, "")
	}()

	c.manager.SubscribeUser(c.ctx, c.userID)

	if c.channelType == ChannelTextTopic {
		c.manager.JoinToTextTopic(c.ctx, c.userID, domain.TopicID(c.channelID))
	}

	for {
		_, rawEvent, err := c.ws.Read(ctx)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				logger.Error(err, "[connection] ws.Read")
			}
			return
		}

		var event Event
		if err := json.Unmarshal(rawEvent, &event); err != nil {
			c.sendError("invalid json format")
			continue
		}

		switch event.Type {
		case EventPong:
			c.handlePong(event)
		case EventJoin:
			c.handleJoin(event)
		case EventLeave:
			c.handleLeave(event)
		default:
			c.sendError("unknown event type")
		}
	}
}

func (c *Connection) WritePump(ctx context.Context) {
	defer c.ws.Close(websocket.StatusNormalClosure, "")

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			c.sendPing()
		case message, ok := <-c.send:
			if !ok {
				return
			}

			bytes, err := render.JsonBytes(message)
			if err != nil {
				logger.Error(err, "[connection] render.JsonBytes")
				continue
			}

			if err := c.ws.Write(ctx, websocket.MessageText, bytes); err != nil {
				logger.Error(err, "[connection] ws.Write")
				return
			}
		}
	}
}

func (c *Connection) sendPing() {
	event := Event{
		Type: EventPing,
	}
	bytes, err := render.JsonBytes(event)
	if err != nil {
		logger.Error(err, "[connection] render.JsonBytes ping")
		return
	}
	c.ws.Write(context.Background(), websocket.MessageText, bytes)
}

func (c *Connection) handlePong(event Event) {
	if c.channelType == ChannelTextTopic {
		topicID := domain.TopicID(event.ChannelID)
		if err := c.manager.JoinToTextTopic(c.ctx, c.userID, topicID); err != nil {
			logger.Error(err, "[connection] handlePong: JoinToTextTopic")
		}
	}
}

func (c *Connection) handleJoin(event Event) {
	if c.channelType == ChannelTextTopic {
		topicID := domain.TopicID(event.ChannelID)
		if err := c.manager.JoinToTextTopic(c.ctx, c.userID, topicID); err != nil {
			logger.Error(err, "[connection] handleJoin: JoinToTextTopic")
		}
	}
}

func (c *Connection) handleLeave(event Event) {
	if c.channelType == ChannelTextTopic {
		topicID := domain.TopicID(event.ChannelID)
		if err := c.manager.LeaveFromTextTopic(c.ctx, c.userID, topicID); err != nil {
			logger.Error(err, "[connection] handleLeave: LeaveFromTextTopic")
		}
	}
}

func (c *Connection) sendError(msg string) {
	event := Event{
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

func (c *Connection) IsActive() bool {
	return c.active
}
