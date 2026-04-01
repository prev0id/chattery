package signaling_api

type EventType string

const (
	EventJoin    EventType = "join"
	EventLeave   EventType = "leave"
	EventMessage EventType = "message"
	EventError   EventType = "error"
)

type ChannelType string

const (
	ChannelDM     ChannelType = "dm"
	ChannelServer ChannelType = "server"
)

type WSEvent struct {
	Type        EventType   `json:"type"`
	ChannelType ChannelType `json:"channel_type,omitempty"`
	ChannelID   string      `json:"channel_id,omitempty"`
	Message     any         `json:"message,omitempty"`
	Error       string      `json:"error,omitempty"`
}
