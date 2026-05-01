package voice_topic

import (
	"encoding/json"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
)

type nodeEventType string

const (
	nodeEventJoin   nodeEventType = "join"
	nodeEventLeave  nodeEventType = "leave"
	nodeEventSignal nodeEventType = "signal"
)

type nodeEvent struct {
	Type      nodeEventType   `json:"type"`
	EventType event_desc.Type `json:"event_type,omitzero"`
	Payload   json.RawMessage `json:"payload,omitzero"`
	UserID    domain.UserID   `json:"user_id"`
	TopicID   domain.TopicID  `json:"topic_id"`
}
