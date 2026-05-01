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
	GetString(ctx context.Context, key string) (string, error)
	SetI64(ctx context.Context, key string, value int64, expiration time.Duration) error
	SetNXString(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
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
	defer close(events)

	go func() {
		r.client.Subscribe(ctx, userChannelKey(userID), sink, done)
	}()

	for {
		select {
		case <-ctx.Done():
			close(done)
			return

		case rawMessage, ok := <-sink:
			if !ok {
				return
			}

			event, err := bind.JSONString[event_desc.Event](rawMessage)
			if err != nil {
				logger.Error(err, "[redis_adapter] bind.JSONString event_desc.Event")
				continue
			}

			select {
			case <-ctx.Done():
				close(done)
				return
			case events <- event:
			}
		}
	}
}

func (r *Adapter) PublishToUser(ctx context.Context, userID domain.UserID, event []byte) error {
	if err := r.client.Publish(ctx, userChannelKey(userID), string(event)); err != nil {
		return errutil.E(err).Debug("r.client.Publish")
	}

	return nil
}

func (r *Adapter) ClaimVoiceTopicOwner(ctx context.Context, topicID domain.TopicID, nodeID string, expiration time.Duration) (bool, error) {
	wasSet, err := r.client.SetNXString(ctx, voiceTopicOwnerKey(topicID), nodeID, expiration)
	if err != nil {
		return false, errutil.E(err).Debug("r.client.SetNXString")
	}
	return wasSet, nil
}

func (r *Adapter) GetVoiceTopicOwner(ctx context.Context, topicID domain.TopicID) (string, error) {
	nodeID, err := r.client.GetString(ctx, voiceTopicOwnerKey(topicID))
	if err != nil {
		return "", errutil.E(err).Debug("r.client.GetString")
	}
	return nodeID, nil
}

func (r *Adapter) RefreshVoiceTopicOwner(ctx context.Context, topicID domain.TopicID, expiration time.Duration) error {
	if err := r.client.Expire(ctx, voiceTopicOwnerKey(topicID), expiration); err != nil {
		return errutil.E(err).Debug("r.client.Expire")
	}
	return nil
}

func (r *Adapter) DeleteVoiceTopicOwner(ctx context.Context, topicID domain.TopicID) error {
	if err := r.client.Delete(ctx, voiceTopicOwnerKey(topicID)); err != nil {
		return errutil.E(err).Debug("r.client.Delete")
	}
	return nil
}

func (r *Adapter) AddVoiceParticipant(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	if err := r.client.ZAddI64(ctx, voiceTopicParticipantsKey(topicID), userID.I64()); err != nil {
		return errutil.E(err).Debug("r.client.ZAddI64")
	}
	return nil
}

func (r *Adapter) RemoveVoiceParticipant(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	if err := r.client.ZRemoveI64(ctx, voiceTopicParticipantsKey(topicID), userID.I64()); err != nil {
		return errutil.E(err).Debug("r.client.ZRemoveI64")
	}
	return nil
}

func (r *Adapter) ListVoiceParticipants(ctx context.Context, topicID domain.TopicID, threshold time.Duration) ([]domain.UserID, error) {
	participants, err := r.client.ZMembersI64(ctx, voiceTopicParticipantsKey(topicID), threshold)
	if err != nil {
		return nil, errutil.E(err).Debug("r.client.ZMembersI64")
	}

	result := make([]domain.UserID, 0, len(participants))
	for _, participant := range participants {
		result = append(result, domain.UserID(participant))
	}
	return result, nil
}

func (r *Adapter) PublishToVoiceNode(ctx context.Context, nodeID string, event []byte) error {
	if err := r.client.Publish(ctx, voiceNodeChannelKey(nodeID), string(event)); err != nil {
		return errutil.E(err).Debug("r.client.Publish")
	}
	return nil
}

func (r *Adapter) SubscribeToVoiceNode(ctx context.Context, nodeID string, events chan<- string) {
	sink := make(chan string)
	done := make(chan struct{})
	defer close(events)

	go func() {
		r.client.Subscribe(ctx, voiceNodeChannelKey(nodeID), sink, done)
	}()

	for {
		select {
		case <-ctx.Done():
			close(done)
			return
		case rawMessage, ok := <-sink:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				close(done)
				return
			case events <- rawMessage:
			}
		}
	}
}

func sessionKey(session domain.Session) string {
	return "Session_" + session.String()
}

func userChannelKey(userID domain.UserID) string {
	return "user:" + strconv.FormatInt(userID.I64(), 10)
}

func voiceTopicOwnerKey(topicID domain.TopicID) string {
	return "voice:topic:" + topicID.String() + ":owner"
}

func voiceTopicParticipantsKey(topicID domain.TopicID) string {
	return "voice:topic:" + topicID.String() + ":participants"
}

func voiceNodeChannelKey(nodeID string) string {
	return "voice:node:" + nodeID + ":events"
}
