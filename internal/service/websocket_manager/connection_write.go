package websocket_manager

import (
	"context"
	"time"

	"github.com/coder/websocket"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

func (c *Connection) WritePump() {
	pingTicker := time.NewTicker(pingInterval)

	defer func() {
		pingTicker.Stop()
		c.cancel()
		if err := c.ws.Close(websocket.StatusNormalClosure, ""); err != nil {
			logger.Error(err, "[connection] WritePump: ws.Close")
		}
	}()

	for c.writeEvent(pingTicker) {
	}
}

func (c *Connection) SendEvent(ctx context.Context, event *event_desc.Event) {
	if event == nil {
		return
	}

	select {
	case <-ctx.Done():
	case <-c.ctx.Done():
	case c.toSend <- event:
	}
}

func (c *Connection) writeToWebsocket(ctx context.Context, event *event_desc.Event) {
	bytes, err := render.JSONBytes(event)
	if err != nil {
		logger.Error(err, "[connection] render.JSONBytes")
		return
	}

	if err := c.ws.Write(ctx, websocket.MessageText, bytes); err != nil {
		logger.Error(err, "[connection] ws.Write")
	}
}

func (c *Connection) writeEvent(ping *time.Ticker) bool {
	select {
	case <-c.ctx.Done():
		return false
	case <-ping.C:
		c.channelMutex.RLock()
		sinceLastPong := time.Since(c.lastPong)
		c.channelMutex.RUnlock()

		if sinceLastPong > pongTimeout {
			return false
		}
		c.sendPing(c.ctx)
	case event, ok := <-c.toSend:
		if !ok {
			return false
		}
		c.writeToWebsocket(c.ctx, event)
	}
	return true
}

func (c *Connection) sendPing(ctx context.Context) {
	ping := &event_desc.Event{Type: event_desc.TypePing}
	c.writeToWebsocket(ctx, ping)
}

func (c *Connection) sendError(message string) {
	payload := event_desc.ErrorPayload{Message: message}

	renderedPayload, err := render.JSONBytes(payload)
	if err != nil {
		logger.Error(err, "[connection] render.JSONBytes error payload")
		return
	}

	event := &event_desc.Event{
		Type:    event_desc.TypeError,
		Payload: renderedPayload,
	}
	c.SendEvent(c.ctx, event)
}
