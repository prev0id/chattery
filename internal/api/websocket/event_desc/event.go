package event_desc

import "encoding/json"

type Type string

const (
	TypePing    Type = "ping"
	TypePong    Type = "pong"
	TypeError   Type = "error"
	TypeJoin    Type = "join"
	TypeLeave   Type = "leave"
	TypeMessage Type = "message"
)

type Event struct {
	Type    Type            `json:"type"`
	Channel Channel         `json:"channel,omitzero"`
	Payload json.RawMessage `json:"payload,omitzero"`
}

type MessagePayload struct {
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Sender    UserInfo `json:"sender"`
	ID        int64    `json:"id"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type UserInfo struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	ID       int64  `json:"id"`
}
