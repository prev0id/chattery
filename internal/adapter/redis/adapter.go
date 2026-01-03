package redisadapter

import (
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
	"context"
	"strconv"
	"time"
)

type client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Publish(ctx context.Context, channel string, message string) error
	Subscribe(ctx context.Context, channel string, sink chan<- string)
}

type Adapter struct {
	client client
}

func NewRedisAdapter(client client) *Adapter {
	return &Adapter{client: client}
}

func (r *Adapter) WriteSession(ctx context.Context, session domain.Session, user domain.Username, expiration time.Duration) error {
	if err := r.client.Set(ctx, sessionKey(session), user.String(), expiration); err != nil {
		return errors.E(err).Debug("r.client.Set")
	}
	return nil
}

func (r *Adapter) UsernameForSession(ctx context.Context, session domain.Session) (domain.Username, error) {
	user, err := r.client.Get(ctx, sessionKey(session))
	if err != nil {
		return domain.UserUnknown, errors.E(err).Debug("r.client.Get")
	}
	return domain.Username(user), nil
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
