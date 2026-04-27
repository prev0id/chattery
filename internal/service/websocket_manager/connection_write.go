package websocket_manager

import (
	"context"
	"time"

	"github.com/coder/websocket"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

func (c *Connection) WritePump(ctx context.Context) {
	pingTicker := time.NewTicker(pingInterval)

	defer func() {
		pingTicker.Stop()
		if err := c.ws.Close(websocket.StatusNormalClosure, ""); err != nil {
			logger.Error(err, "[connection] WritePump: ws.Close")
		}
	}()

	for c.writeEvent(ctx, pingTicker) {
	}
}

func (c *Connection) SendEvent(ctx context.Context, event *event_desc.Event) {
	bytes, err := render.JSONBytes(event)
	if err != nil {
		logger.Error(err, "[connection] render.JSONBytes")
		return
	}

	if err := c.ws.Write(ctx, websocket.MessageText, bytes); err != nil {
		logger.Error(err, "[connection] ws.Write")
	}
}

func (c *Connection) writeEvent(ctx context.Context, ping *time.Ticker) bool {
	select {
	case <-ctx.Done():
		return false
	case <-ping.C:
		c.channelMutex.RLock()
		sinceLastPong := time.Since(c.lastPong)
		c.channelMutex.RUnlock()

		if sinceLastPong > pongTimeout {
			return false
		}
		c.sendPing(ctx)
	case event, ok := <-c.toSend:
		if !ok {
			return false
		}
		c.SendEvent(ctx, event)
	}
	return true
}

func (c *Connection) sendPing(ctx context.Context) {
	ping := &event_desc.Event{Type: event_desc.TypePing}
	c.SendEvent(ctx, ping)
}

func (c *Connection) sendError(message string) {
	payload := event_desc.ErrorPayload{Message: message}

	event, err := render.Event(event_desc.TypeError, event_desc.Channel{}, payload)
	if err != nil {
		logger.Error(err, "[connection] render.JSONBytes error")
		return
	}

	if err := c.ws.Write(context.Background(), websocket.MessageText, event); err != nil {
		logger.Error(err, "[connection] sendError: ws.Write")
	}
}
