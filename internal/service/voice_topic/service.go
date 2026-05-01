package voice_topic

import (
	"context"
	"encoding/json"
	"maps"
	"sync"
	"time"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

type redis interface {
	ClaimVoiceTopicOwner(ctx context.Context, topicID domain.TopicID, nodeID string, expiration time.Duration) (bool, error)
	GetVoiceTopicOwner(ctx context.Context, topicID domain.TopicID) (string, error)
	RefreshVoiceTopicOwner(ctx context.Context, topicID domain.TopicID, expiration time.Duration) error
	DeleteVoiceTopicOwner(ctx context.Context, topicID domain.TopicID) error
	AddVoiceParticipant(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
	RemoveVoiceParticipant(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
	ListVoiceParticipants(ctx context.Context, topicID domain.TopicID, threshold time.Duration) ([]domain.UserID, error)
	PublishToVoiceNode(ctx context.Context, nodeID string, event []byte) error
	SubscribeToVoiceNode(ctx context.Context, nodeID string, events chan<- string)
	PublishToUser(ctx context.Context, userID domain.UserID, event []byte) error
}

type serverService interface {
	ValidateAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID, topicType domain.TopicType) error
}

type userCache interface {
	GetByID(userID domain.UserID) (*domain.User, error)
}

type Service struct {
	topics   map[domain.TopicID]*topic
	redis    redis
	server   serverService
	user     userCache
	webrtc   *webrtcAPI
	nodeID   string
	ownerTTL time.Duration
	m        sync.RWMutex
}

func New(redisAdapter redis, server serverService, user userCache, cfg *config.Config) *Service {
	ownerTTL := cfg.Voice.OwnerTTL
	if ownerTTL <= 0 {
		ownerTTL = 30 * time.Second
	}

	return &Service{
		redis:    redisAdapter,
		server:   server,
		user:     user,
		nodeID:   cfg.Voice.NodeID,
		ownerTTL: ownerTTL,
		webrtc:   newWebRTCAPI(cfg),
		topics:   make(map[domain.TopicID]*topic),
	}
}

func (s *Service) Start(ctx context.Context) {
	go s.subscribeNodeEvents(ctx)
	go s.refreshOwnedTopics(ctx)
}

func (s *Service) Join(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := s.server.ValidateAccessToTopic(ctx, userID, topicID, domain.TopicTypeVoice); err != nil {
		return err
	}

	owner, err := s.ownerNode(ctx, topicID)
	if err != nil {
		return err
	}

	if owner != s.nodeID {
		return s.publishNodeEvent(ctx, owner, nodeEvent{
			Type:    nodeEventJoin,
			TopicID: topicID,
			UserID:  userID,
		})
	}

	return s.joinLocal(ctx, topicID, userID)
}

func (s *Service) Leave(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	owner, err := s.redis.GetVoiceTopicOwner(ctx, topicID)
	if errutil.Is(errutil.NotFound, err) {
		return nil
	}
	if err != nil {
		return errutil.E(err).Debug("s.redis.GetVoiceTopicOwner")
	}

	if owner != s.nodeID {
		return s.publishNodeEvent(ctx, owner, nodeEvent{
			Type:    nodeEventLeave,
			TopicID: topicID,
			UserID:  userID,
		})
	}

	return s.leaveLocal(ctx, topicID, userID)
}

func (s *Service) HandleSignal(ctx context.Context, userID domain.UserID, event *event_desc.Event) error {
	topicID := domain.TopicID(event.Channel.ID)
	if topicID == 0 {
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Message("voice event channel is required")
	}

	if err := s.server.ValidateAccessToTopic(ctx, userID, topicID, domain.TopicTypeVoice); err != nil {
		return err
	}

	owner, err := s.ownerNode(ctx, topicID)
	if err != nil {
		return err
	}

	if owner != s.nodeID {
		return s.publishNodeEvent(ctx, owner, nodeEvent{
			Type:      nodeEventSignal,
			TopicID:   topicID,
			UserID:    userID,
			EventType: event.Type,
			Payload:   event.Payload,
		})
	}

	return s.handleSignalLocal(ctx, topicID, userID, event.Type, event.Payload)
}

func (s *Service) ownerNode(ctx context.Context, topicID domain.TopicID) (string, error) {
	owner, err := s.redis.GetVoiceTopicOwner(ctx, topicID)
	if err == nil {
		return owner, nil
	}
	if !errutil.Is(errutil.NotFound, err) {
		return "", errutil.E(err).Debug("s.redis.GetVoiceTopicOwner")
	}

	claimed, err := s.redis.ClaimVoiceTopicOwner(ctx, topicID, s.nodeID, s.ownerTTL)
	if err != nil {
		return "", errutil.E(err).Debug("s.redis.ClaimVoiceTopicOwner")
	}
	if claimed {
		return s.nodeID, nil
	}

	owner, err = s.redis.GetVoiceTopicOwner(ctx, topicID)
	if err != nil {
		return "", errutil.E(err).Debug("s.redis.GetVoiceTopicOwner", "after claim race")
	}
	return owner, nil
}

func (s *Service) joinLocal(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	topic := s.getTopic(topicID)
	joined := topic.addParticipant(userID)
	if joined {
		if err := s.redis.AddVoiceParticipant(ctx, topicID, userID); err != nil {
			return errutil.E(err).Debug("s.redis.AddVoiceParticipant")
		}
	}

	if err := s.broadcastParticipant(ctx, topicID, userID, event_desc.TypeVoiceJoined); err != nil {
		return err
	}
	if joined {
		s.sendExistingParticipants(ctx, topicID, userID)
	}
	return nil
}

func (s *Service) leaveLocal(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error {
	s.m.RLock()
	topic := s.topics[topicID]
	s.m.RUnlock()
	if topic == nil {
		return nil
	}

	if topic.removeParticipant(userID) {
		if err := s.redis.RemoveVoiceParticipant(ctx, topicID, userID); err != nil {
			logger.ErrorCtx(ctx, err, "s.redis.RemoveVoiceParticipant")
		}
		if err := s.broadcastParticipant(ctx, topicID, userID, event_desc.TypeVoiceLeft); err != nil {
			return err
		}
	}

	if topic.empty() {
		s.m.Lock()
		delete(s.topics, topicID)
		s.m.Unlock()
		if err := s.redis.DeleteVoiceTopicOwner(ctx, topicID); err != nil {
			logger.ErrorCtx(ctx, err, "s.redis.DeleteVoiceTopicOwner")
		}
	}

	return nil
}

func (s *Service) handleSignalLocal(ctx context.Context, topicID domain.TopicID, userID domain.UserID, eventType event_desc.Type, payload json.RawMessage) error {
	if !s.hasParticipant(topicID, userID) {
		if err := s.joinLocal(ctx, topicID, userID); err != nil {
			return err
		}
	}

	topic := s.getTopic(topicID)

	switch eventType {
	case event_desc.TypeVoiceOffer:
		return topic.handleOffer(ctx, userID, payload)
	case event_desc.TypeVoiceAnswer:
		return topic.handleAnswer(ctx, userID, payload)
	case event_desc.TypeVoiceICECandidate:
		return topic.handleICECandidate(ctx, userID, payload)
	default:
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Message("unknown voice signal event type")
	}
}

func (s *Service) broadcastParticipant(ctx context.Context, topicID domain.TopicID, userID domain.UserID, eventType event_desc.Type) error {
	payload := s.convertParticipantPayload(topicID, userID)
	event, err := render.Event(eventType, event_desc.Channel{
		Type: event_desc.ChannelVoiceTopic,
		ID:   topicID.I64(),
	}, payload)
	if err != nil {
		return errutil.E(err).Debug("render.Event")
	}

	participants, err := s.redis.ListVoiceParticipants(ctx, topicID, s.ownerTTL)
	if err != nil {
		return errutil.E(err).Debug("s.redis.ListVoiceParticipants")
	}

	for _, participant := range participants {
		if err := s.redis.PublishToUser(ctx, participant, event); err != nil {
			logger.ErrorCtx(ctx, err, "s.redis.PublishToUser")
		}
	}
	return nil
}

func (s *Service) sendExistingParticipants(ctx context.Context, topicID domain.TopicID, to domain.UserID) {
	participants, err := s.redis.ListVoiceParticipants(ctx, topicID, s.ownerTTL)
	if err != nil {
		logger.ErrorCtx(ctx, err, "s.redis.ListVoiceParticipants")
		return
	}

	for _, participant := range participants {
		if participant == to {
			continue
		}

		payload := s.convertParticipantPayload(topicID, participant)
		event, err := render.Event(event_desc.TypeVoiceJoined, event_desc.Channel{
			Type: event_desc.ChannelVoiceTopic,
			ID:   topicID.I64(),
		}, payload)
		if err != nil {
			logger.ErrorCtx(ctx, err, "render.Event")
			continue
		}
		if err := s.redis.PublishToUser(ctx, to, event); err != nil {
			logger.ErrorCtx(ctx, err, "s.redis.PublishToUser")
		}
	}
}

func (s *Service) convertParticipantPayload(topicID domain.TopicID, userID domain.UserID) event_desc.VoiceParticipantPayload {
	return event_desc.VoiceParticipantPayload{
		UserID:  userID.I64(),
		TopicID: topicID.I64(),
		NodeID:  s.nodeID,
		Sender:  s.getUserInfo(userID),
	}
}

func (s *Service) getUserInfo(userID domain.UserID) event_desc.UserInfo {
	user, err := s.user.GetByID(userID)
	if err != nil {
		logger.Error(err, "s.user.GetByID")
		return event_desc.UserInfo{ID: userID.I64()}
	}

	return event_desc.UserInfo{
		Username: user.Username.String(),
		Avatar:   "/v1/image/" + user.Username.String() + ".png",
		ID:       user.ID.I64(),
	}
}

func (s *Service) getTopic(topicID domain.TopicID) *topic {
	s.m.Lock()
	defer s.m.Unlock()

	if topic, ok := s.topics[topicID]; ok {
		return topic
	}

	topic := newTopic(topicID, s)
	s.topics[topicID] = topic
	return topic
}

func (s *Service) hasParticipant(topicID domain.TopicID, userID domain.UserID) bool {
	s.m.RLock()
	topic := s.topics[topicID]
	s.m.RUnlock()
	return topic != nil && topic.hasParticipant(userID)
}

func (s *Service) publishNodeEvent(ctx context.Context, owner string, event nodeEvent) error {
	payload, err := render.JSONBytes(event)
	if err != nil {
		return errutil.E(err).Debug("render.JSONBytes")
	}
	if err := s.redis.PublishToVoiceNode(ctx, owner, payload); err != nil {
		return errutil.E(err).Debug("s.redis.PublishToVoiceNode")
	}
	return nil
}

func (s *Service) subscribeNodeEvents(ctx context.Context) {
	events := make(chan string)
	go s.redis.SubscribeToVoiceNode(ctx, s.nodeID, events)

	for rawEvent := range events {
		event, err := bind.JSONString[nodeEvent](rawEvent)
		if err != nil {
			logger.ErrorCtx(ctx, err, "[voice_topic] bind.JSONString nodeEvent")
			continue
		}
		if err := s.handleNodeEvent(ctx, event); err != nil {
			logger.ErrorCtx(ctx, err, "[voice_topic] handleNodeEvent")
		}
	}
}

func (s *Service) handleNodeEvent(ctx context.Context, event *nodeEvent) error {
	switch event.Type {
	case nodeEventJoin:
		return s.joinLocal(ctx, event.TopicID, event.UserID)
	case nodeEventLeave:
		return s.leaveLocal(ctx, event.TopicID, event.UserID)
	case nodeEventSignal:
		return s.handleSignalLocal(ctx, event.TopicID, event.UserID, event.EventType, event.Payload)
	default:
		return errutil.E().
			Kind(errutil.InvalidRequest).
			Message("unknown voice node event type")
	}
}

func (s *Service) refreshOwnedTopics(ctx context.Context) {
	if s.ownerTTL <= 0 {
		return
	}

	ticker := time.NewTicker(s.ownerTTL / 3)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.refreshOwnedTopicsOnce(ctx)
		}
	}
}

func (s *Service) refreshOwnedTopicsOnce(ctx context.Context) {
	s.m.RLock()
	topics := make(map[domain.TopicID]*topic, len(s.topics))
	maps.Copy(topics, s.topics)
	s.m.RUnlock()

	for topicID, topic := range topics {
		if err := s.redis.RefreshVoiceTopicOwner(ctx, topicID, s.ownerTTL); err != nil {
			logger.ErrorCtx(ctx, err, "s.redis.RefreshVoiceTopicOwner")
		}
		for _, participant := range topic.participantIDs() {
			if err := s.redis.AddVoiceParticipant(ctx, topicID, participant); err != nil {
				logger.ErrorCtx(ctx, err, "s.redis.AddVoiceParticipant")
			}
		}
	}
}
