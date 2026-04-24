package websocket_manager

import (
	"context"
	"sync"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

type dmService interface {
	GetParticipants(ctx context.Context, dmID domain.DMID) ([]domain.UserID, error)
}

type serverService interface {
	UserHasAccessToTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
}

type textTopicService interface {
	AddUserInTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	RemoveUserFromTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
}

type userCache interface {
	GetByID(userID domain.UserID) (*domain.User, error)
}

type redis interface {
	SubscribeToUser(ctx context.Context, userID domain.UserID, dst chan<- *domain.UserMessage)
	PublishToUser(ctx context.Context, userID domain.UserID, message *domain.UserMessage) error
	RemoveUserFromTextTopic(ctx context.Context, topicID domain.TopicID, userID domain.UserID) error
}

type WebsocketManager struct {
	redis            redis
	dmService        dmService
	serverService    serverService
	textTopicService textTopicService
	user             userCache

	connections     map[domain.UserID]map[*Connection]bool
	userSubscribers map[domain.UserID]chan *domain.UserMessage

	m sync.RWMutex
}

func New(redisAdapter redis, dmSvc dmService, serverSvc serverService, textTopicSvc textTopicService, userCache userCache) *WebsocketManager {
	return &WebsocketManager{
		redis:            redisAdapter,
		dmService:        dmSvc,
		serverService:    serverSvc,
		textTopicService: textTopicSvc,
		user:             userCache,
		connections:      make(map[domain.UserID]map[*Connection]bool),
		userSubscribers:  make(map[domain.UserID]chan *domain.UserMessage),
	}
}

func (m *WebsocketManager) Redis() redis {
	return m.redis
}

func (m *WebsocketManager) UserHasAccessToDM(ctx context.Context, userID domain.UserID, dmID domain.DMID) error {
	participants, err := m.dmService.GetParticipants(ctx, dmID)
	if err != nil {
		return errors.E(err).Debug("m.dmService.GetParticipants")
	}

	for _, p := range participants {
		if p == userID {
			return nil
		}
	}

	return errors.E().Kind(errors.InvalidRequest).Debug("user is not participant of dm")
}

func (m *WebsocketManager) UserHasAccessToTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := m.serverService.UserHasAccessToTopic(ctx, userID, topicID); err != nil {
		return errors.E(err).Debug("m.serverService.UserHasAccessToTopic")
	}
	return nil
}

func (m *WebsocketManager) JoinToTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	if err := m.textTopicService.AddUserInTextTopic(ctx, userID, topicID); err != nil {
		return errors.E(err).Debug("m.textTopicService.AddUserInTextTopic")
	}
	return nil
}

func (m *WebsocketManager) LeaveFromTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error {
	return m.textTopicService.RemoveUserFromTextTopic(ctx, userID, topicID)
}

func (m *WebsocketManager) RegisterConnection(conn *Connection) {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.connections[conn.userID]; !ok {
		m.connections[conn.userID] = make(map[*Connection]bool)
	}
	m.connections[conn.userID][conn] = true
}

func (m *WebsocketManager) UnregisterConnection(conn *Connection) {
	m.m.Lock()
	defer m.m.Unlock()

	if conns, ok := m.connections[conn.userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(m.connections, conn.userID)
			m.stopUserSubscriber(conn.userID)
		}
	}
}

func (m *WebsocketManager) SubscribeUser(ctx context.Context, userID domain.UserID) {
	m.m.Lock()
	defer m.m.Unlock()

	if _, ok := m.userSubscribers[userID]; ok {
		return
	}

	dst := make(chan *domain.UserMessage)
	m.userSubscribers[userID] = dst

	go m.redis.SubscribeToUser(ctx, userID, dst)
	go m.handleUserMessages(userID, dst)
}

func (m *WebsocketManager) stopUserSubscriber(userID domain.UserID) {
	if ch, ok := m.userSubscribers[userID]; ok {
		close(ch)
		delete(m.userSubscribers, userID)
	}
}

func (m *WebsocketManager) handleUserMessages(userID domain.UserID, src chan *domain.UserMessage) {
	for msg := range src {
		m.m.RLock()
		conns := m.connections[userID]
		m.m.RUnlock()

		for conn := range conns {
			event := m.userMessageToEvent(msg, conn.channelType, conn.channelID)
			if event != nil {
				select {
				case conn.send <- event:
				default:
				}
			}
		}
	}
}

func (m *WebsocketManager) userMessageToEvent(msg *domain.UserMessage, connChannelType ChannelType, connChannelID int64) *Event {
	switch msg.Type {
	case domain.UserMessageTypeDM:
		if connChannelType != ChannelDM {
			return nil
		}
		if int64(msg.DMMessage.DMID.I64()) != connChannelID {
			return nil
		}
		return m.dmToEventMessage(msg.DMMessage)
	case domain.UserMessageTypeTopic:
		if connChannelType != ChannelTextTopic {
			return nil
		}
		if int64(msg.TopicMsg.TopicID.I64()) != connChannelID {
			return nil
		}
		return m.topicToEventMessage(msg.TopicMsg)
	default:
		return nil
	}
}

func (m *WebsocketManager) dmToEventMessage(msg *domain.DMMessage) *Event {
	sender := UserInfo{}
	if user, err := m.user.GetByID(msg.SenderID); err == nil {
		sender = UserInfo{
			ID:       int64(user.ID),
			Username: user.Username.String(),
			Avatar:   user.AvatarID.String(),
		}
	}

	return &Event{
		Type:        EventMessage,
		ChannelType: ChannelDM,
		ChannelID:   msg.DMID.I64(),
		Message: &EventData{
			ID:        msg.ID.I64(),
			Sender:    sender,
			Text:      msg.Text,
			CreatedAt: render.Timestamp(msg.CreatedAt),
		},
	}
}

func (m *WebsocketManager) topicToEventMessage(msg *domain.TopicMessage) *Event {
	sender := UserInfo{}
	if user, err := m.user.GetByID(msg.SenderID); err == nil {
		sender = UserInfo{
			ID:       user.ID.I64(),
			Username: user.Username.String(),
			Avatar:   user.AvatarID.String(),
		}
	}

	return &Event{
		Type:        EventMessage,
		ChannelType: ChannelTextTopic,
		ChannelID:   msg.TopicID.I64(),
		Message: &EventData{
			ID:        msg.ID.I64(),
			Sender:    sender,
			Text:      msg.Text,
			CreatedAt: render.Timestamp(msg.CreatedAt),
		},
	}
}
