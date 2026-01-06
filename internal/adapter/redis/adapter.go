package redisadapter

import (
	"context"
	"strconv"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

type client interface {
	GetExI64(ctx context.Context, key string, expiration time.Duration) (int64, error)
	SetI64(ctx context.Context, key string, value int64, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Publish(ctx context.Context, channel string, message string) error
	Subscribe(ctx context.Context, channel string, sink chan<- string)
}

type Adapter struct {
	client client
}

func NewRedisAdapter(client client) *Adapter {
	return &Adapter{client: client}
}

func (r *Adapter) WriteSession(ctx context.Context, session domain.Session, userID domain.UserID, expiration time.Duration) error {
	if err := r.client.SetI64(ctx, sessionKey(session), userID.I64(), expiration); err != nil {
		return errors.E(err).Debug("r.client.SetI64")
	}
	return nil
}

func (r *Adapter) UserIDFromSession(ctx context.Context, session domain.Session, expiration time.Duration) domain.UserID {
	userID, err := r.client.GetExI64(ctx, sessionKey(session), expiration)
	if err != nil {
		logger.Error(err, "r.client.GetExI64")
		return domain.UserIsUnknown
	}
	return domain.UserID(userID)
}

func (r *Adapter) ClearSession(ctx context.Context, session domain.Session) error {
	if err := r.client.Delete(ctx, sessionKey(session)); err != nil {
		return errors.E(err).Debug("r.client.Delete")
	}
	return nil

}

func (r *Adapter) SendMessage(ctx context.Context, chat domain.ChatID, message domain.Message) error {
	renderedMessage, err := render.JsonString(message)
	if err != nil {
		return errors.E(err).Debug("render.JsonString")
	}
	if err := r.client.Publish(ctx, chatKey(chat), renderedMessage); err != nil {
		return errors.E(err).Debug("r.client.Publish")
	}
	return nil
}

func (r *Adapter) Subscribe(ctx context.Context, chat domain.ChatID, dst chan<- *domain.Message) {
	sink := make(chan string)

	go func() {
		r.client.Subscribe(ctx, chatKey(chat), sink)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case rawMessage := <-sink:
			message, err := bind.JsonString[domain.Message](rawMessage)
			if err != nil {
				logger.Error(err, "[redis_adapter] bind.JsonString")
				continue
			}
			dst <- message
		}
	}
}

func sessionKey(session domain.Session) string {
	return "Session_" + session.String()
}

func chatKey(chatID domain.ChatID) string {
	return "Chat_" + strconv.FormatInt(chatID.I64(), 10)
}
