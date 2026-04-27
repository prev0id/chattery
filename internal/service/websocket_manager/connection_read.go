package websocket_manager

import (
	"context"
	"errors"
	"time"

	"github.com/coder/websocket"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/logger"
)

func (c *Connection) ReadPump(ctx context.Context) {
	defer c.cleanup()

	for c.readEvent(ctx) {
	}
}

func (c *Connection) cleanup() {
	c.cancel()
	c.manager.unregisterConnection(c)
	close(c.toSend)
	if err := c.ws.Close(websocket.StatusNormalClosure, ""); err != nil {
		logger.Error(err, "[connection] cleanup: ws.Close")
	}
}

func (c *Connection) readEvent(ctx context.Context) bool {
	_, rawEvent, err := c.ws.Read(ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Error(err, "[connection] ws.Read")
		}
		return false
	}

	event, err := bind.JSONBytes[event_desc.Event](rawEvent)
	if err != nil {
		logger.Error(err, "[connection] bind.JSONBytes")
		return true
	}

	c.handleEvent(ctx, event)
	return true
}

func (c *Connection) handleEvent(ctx context.Context, event *event_desc.Event) {
	switch event.Type {
	case event_desc.TypePong:
		c.handlePong(ctx, event)
	case event_desc.TypeJoin:
		c.handleJoin(ctx, event)
	case event_desc.TypeLeave:
		c.handleLeave(ctx, event)
	default:
		c.sendError("unknown event type")
	}
}

func (c *Connection) handlePong(_ context.Context, _ *event_desc.Event) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()
	c.lastPong = time.Now()
}

func (c *Connection) handleJoin(ctx context.Context, event *event_desc.Event) {
	channel, err := bind.JSONBytes[event_desc.Channel](event.Payload)
	if err != nil {
		c.sendError("invalid channel payload")
		return
	}

	switch channel.Type {
	case event_desc.ChannelDM:
		dmID := domain.DMID(channel.ID)
		if err := c.manager.dm.ValidateAccess(ctx, c.userID, dmID); err != nil {
			c.sendError("no access to dm")
			return
		}
	case event_desc.ChannelTextTopic:
		topicID := domain.TopicID(channel.ID)
		if err := c.manager.server.ValidateAccessToTopic(ctx, c.userID, topicID, domain.TopicTypeText); err != nil {
			c.sendError("no access to topic")
			return
		}
	default:
		c.sendError("invalid channel type")
		return
	}

	c.channelMutex.Lock()
	c.channel = *channel
	c.channelMutex.Unlock()
}

func (c *Connection) handleLeave(_ context.Context, _ *event_desc.Event) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()
	c.channel = event_desc.Channel{}
}
