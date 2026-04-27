package redis

import (
	"context"
	"strconv"
	"time"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
)

type client interface {
	GetExI64(ctx context.Context, key string, expiration time.Duration) (int64, error)
	SetI64(ctx context.Context, key string, value int64, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Publish(ctx context.Context, channel string, message string) error
	Subscribe(ctx context.Context, channel string, sink chan<- string, done <-chan struct{})
	ZMembersI64(ctx context.Context, key string, threshold time.Duration) ([]int64, error)
	ZAddI64(ctx context.Context, key string, value int64) error
	ZRemoveI64(ctx context.Context, key string, value int64) error
}

type Adapter struct {
	client client
}

func NewRedisAdapter(client client) *Adapter {
	return &Adapter{client: client}
}

func (r *Adapter) WriteSession(ctx context.Context, session domain.Session, userID domain.UserID, expiration time.Duration) error {
	if err := r.client.SetI64(ctx, sessionKey(session), userID.I64(), expiration); err != nil {
		return errutil.E(err).Debug("r.client.SetI64")
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
		return errutil.E(err).Debug("r.client.Delete")
	}
	return nil
}

func (r *Adapter) SubscribeToUser(ctx context.Context, userID domain.UserID, events chan<- *event_desc.Event) {
	sink := make(chan string)
	done := make(chan struct{})

	go func() {
		r.client.Subscribe(ctx, userChannelKey(userID), sink, done)
	}()

	for {
		select {
		case <-ctx.Done():
			close(done)
			return
		case rawMessage := <-sink:
			event, err := bind.JSONString[event_desc.Event](rawMessage)
			if err != nil {
				logger.Error(err, "[redis_adapter] bind.JSONString event_desc.Event")
				continue
			}
			events <- event
		}
	}
}

func (r *Adapter) PublishToUser(ctx context.Context, userID domain.UserID, event []byte) error {
	if err := r.client.Publish(ctx, userChannelKey(userID), string(event)); err != nil {
		return errutil.E(err).Debug("r.client.Publish")
	}

	return nil
}

func sessionKey(session domain.Session) string {
	return "Session_" + session.String()
}

func userChannelKey(userID domain.UserID) string {
	return "user:" + strconv.FormatInt(userID.I64(), 10)
}
