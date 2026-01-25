package subscriber

import (
	"context"

	"github.com/coder/websocket"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

// Subscriber interface for Subscriber

var _ domain.Subscriber = (*Subscriber)(nil)

type chatService interface {
	PostMessage(ctx context.Context, message *domain.Message) error
	StartListeningToChat(ctx context.Context, sub domain.Subscriber, chat domain.ChatID)
	StopListeningToChat(sub domain.Subscriber)
}

type Subscriber struct {
	user    domain.UserID
	session domain.Session
	ws      *websocket.Conn
	chat    chatService
}

func New(ws *websocket.Conn) *Subscriber {
	return &Subscriber{
		ws: ws,
	}
}

func (sub *Subscriber) WithUserID(user domain.UserID) *Subscriber {
	sub.user = user
	return sub
}

func (sub *Subscriber) WithSession(session domain.Session) *Subscriber {
	sub.session = session
	return sub
}

func (sub *Subscriber) WithChatService(chat chatService) *Subscriber {
	sub.chat = chat
	return sub
}

func (sub *Subscriber) GetSession() domain.Session {
	return sub.session
}

func (sub *Subscriber) GetUserID() domain.UserID {
	return sub.user
}

func (sub *Subscriber) WriteMessage(ctx context.Context, message *domain.Message) error {
	bytes, err := render.JsonBytes(convertMessageToEvent(message))
	if err != nil {
		return errors.E(err).Debug("render.JsonBytes")
	}
	if err := sub.ws.Write(ctx, websocket.MessageText, bytes); err != nil {
		return errors.E(err).Debug("[subscriber] sub.ws.Write")
	}
	return nil
}

func (sub *Subscriber) Read(ctx context.Context) {
	for {
		mt, rawEvent, err := sub.ws.Read(ctx)
		if err != nil {
			logger.Error(err, "[subscriber] sub.ws.Read")
			return
		}
		if mt != websocket.MessageText {
			logger.Error(err, "[subscriber] got unsupported message type")
			return
		}

		event, err := bind.JsonBytes[Event](rawEvent)
		if err != nil {
			logger.Error(err, "[subscriber] bind.JsonBytes")
			return
		}

		switch event.Type {
		case EventTypeJoinChat:
			chat := domain.ChatID(event.ChatID)
			sub.chat.StartListeningToChat(ctx, sub, chat)

		case EventTypeLeaveChat:
			sub.chat.StopListeningToChat(sub)

		case EventTypeMessage:
			message := convertEventToMessage(event, sub.user)
			if err := sub.chat.PostMessage(ctx, message); err != nil {
				logger.Error(err, "[subscriber] post new message")
			}
		}
	}
}
