package websocket_manager

import "chattery/internal/domain"

type Event struct {
	Type        EventType          `json:"type"`
	ChannelType domain.ChannelType `json:"channel_type,omitempty"`
	ChannelID   int64              `json:"channel_id,omitempty"`
	Message     *EventData         `json:"message,omitempty"`
	Error       string             `json:"error,omitempty"`
}

type EventData struct {
	ID        int64    `json:"id"`
	Sender    UserInfo `json:"sender"`
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type EventType string

const (
	EventJoin    EventType = "join"
	EventLeave   EventType = "leave"
	EventMessage EventType = "message"
	EventError   EventType = "error"
	EventPing    EventType = "ping"
	EventPong    EventType = "pong"
)

type ChannelKey struct {
	Type domain.ChannelType
	ID   int64
}

func (ck ChannelKey) DMID() domain.DMID {
	return domain.DMID(ck.ID)
}

func (ck ChannelKey) TopicID() domain.TopicID {
	return domain.TopicID(ck.ID)
}
