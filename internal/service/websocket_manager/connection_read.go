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

func (c *Connection) ReadPump() {
	defer c.cleanup()

	for c.readEvent(c.ctx) {
	}
}

func (c *Connection) cleanup() {
	c.leaveVoiceOnCleanup()
	c.cancel()
	c.manager.unregisterConnection(c)
	if err := c.ws.Close(websocket.StatusNormalClosure, ""); err != nil {
		logger.Error(err, "[connection] cleanup: ws.Close")
	}
}

func (c *Connection) leaveVoiceOnCleanup() {
	c.channelMutex.RLock()
	channel := c.channel
	c.channelMutex.RUnlock()

	if channel.Type != event_desc.ChannelVoiceTopic {
		return
	}
	if err := c.manager.voice.Leave(c.ctx, c.userID, domain.TopicID(channel.ID)); err != nil {
		logger.ErrorCtx(c.ctx, err, "[connection] voice.Leave")
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
		c.handlePong()
	case event_desc.TypeJoin:
		c.handleJoin(ctx, event)
	case event_desc.TypeLeave:
		c.handleLeave(ctx)
	case event_desc.TypeVoiceOffer, event_desc.TypeVoiceAnswer, event_desc.TypeVoiceICECandidate, event_desc.TypeVoiceICECandidates:
		c.handleVoiceSignal(ctx, event)
	default:
		c.sendError("unknown event type")
	}
}

func (c *Connection) handlePong() {
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
	case event_desc.ChannelVoiceTopic:
		topicID := domain.TopicID(channel.ID)
		if c.manager.voice == nil {
			c.sendError("voice service is unavailable")
			return
		}
		if err := c.manager.voice.Join(ctx, c.userID, topicID); err != nil {
			c.sendError(err.Error())
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

func (c *Connection) handleLeave(ctx context.Context) {
	c.channelMutex.RLock()
	channel := c.channel
	c.channelMutex.RUnlock()

	if channel.Type == event_desc.ChannelVoiceTopic && c.manager.voice != nil {
		if err := c.manager.voice.Leave(ctx, c.userID, domain.TopicID(channel.ID)); err != nil {
			c.sendError(err.Error())
		}
	}

	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()
	c.channel = event_desc.Channel{}
}

func (c *Connection) handleVoiceSignal(ctx context.Context, event *event_desc.Event) {
	c.channelMutex.RLock()
	channel := c.channel
	c.channelMutex.RUnlock()

	if event.Channel.ID == 0 && channel.Type == event_desc.ChannelVoiceTopic {
		event.Channel = channel
	}
	if event.Channel.Type != event_desc.ChannelVoiceTopic {
		c.sendError("voice signal requires voice topic channel")
		return
	}
	if err := c.manager.voice.HandleSignal(ctx, c.userID, event); err != nil {
		c.sendError(err.Error())
	}
}
