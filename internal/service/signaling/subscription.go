package signaling

import (
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
	"context"

	"github.com/coder/websocket"
)

type Subscriber struct {
	user   domain.Username
	ws     *websocket.Conn
	cancel func()
	reads  chan *domain.Event
	writes chan *domain.Event
}

func (sub *Subscriber) Writer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			err := errors.E(ctx.Err()).User(sub.user)
			logger.Error(err, "[subscriber] ctx.Done")
			return

		case event := <-sub.writes:
			message, err := render.JsonBytes(event)
			if err != nil {
				err = errors.E(err).User(sub.user)
				logger.Error(err, "[subscriber] render.JsonBytes")
				continue
			}
			if err := sub.ws.Write(ctx, websocket.MessageText, message); err != nil {
				err = errors.E(err).User(sub.user)
				logger.Error(err, "[subscriber] sub.ws.Write")
				sub.cancel()
				return
			}
		}
	}
}

func (sub *Subscriber) Reader(ctx context.Context) {
	defer sub.cancel()

	for {
		mt, rawMessage, err := sub.ws.Read(ctx)
		if err != nil {
			logger.Error(err, "sub.ws.Read")
			return
		}
		if mt != websocket.MessageText {
			logger.Error(err, "got unsupported message type")
			return
		}

		event, err := bind.JsonBytes[domain.Event](rawMessage)
		if err != nil {
			logger.Error(err, "bind.JsonBytes")
			return
		}

		sub.reads <- event
	}
}
