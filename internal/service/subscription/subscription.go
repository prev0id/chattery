package subscription

import (
	"context"

	"github.com/coder/websocket"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

type callback func(ctx context.Context, event *domain.Event) error

type Subscriber struct {
	user    domain.UserID
	session domain.Session
	ws      *websocket.Conn
	cancel  func()
	reads   chan *domain.Event

	eventListeners map[domain.EventType]callback
}

func New(ctx context.Context, userID domain.UserID, session domain.Session, ws *websocket.Conn) (context.Context, *Subscriber) {
	ctx, cancel := context.WithCancel(ctx)

	sub := &Subscriber{
		user:           userID,
		ws:             ws,
		cancel:         cancel,
		eventListeners: make(map[domain.EventType]callback),
	}
	return ctx, sub
}

func (sub *Subscriber) Write(ctx context.Context, event domain.Event) error {
	message, err := render.JsonBytes(event)
	if err != nil {
		return errors.E(err).Debug("render.JsonBytes")
	}
	if err := sub.ws.Write(ctx, websocket.MessageText, message); err != nil {
		sub.cancel()
		return errors.E(err).Debug("[subscriber] sub.ws.Write")
	}
	return nil
}

func (sub *Subscriber) SubscribeToEvent(ctx context.Context, type_ domain.EventType, callback callback) {
	sub.eventListeners[type_] = callback
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
