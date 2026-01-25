package domain

import "time"

type ChatID int64

const ChatUnknown ChatID = 0

func (id ChatID) I64() int64 { return int64(id) }

type ChatType string

func (t ChatType) String() string { return string(t) }

const (
	ChatTypePrivate ChatType = "private"
	ChatTypePublic  ChatType = "public"
)

type Chat struct {
	ID   ChatID
	Name string
	Type ChatType
}

type Participant struct {
	UserID UserID
	Chat   ChatID
	Role   ChatRole
}

type ChatRole string

const (
	ChatRoleOwner       ChatRole = "owner"
	ChatRoleParticipant ChatRole = "participant"
)

func (role ChatRole) String() string { return string(role) }

type MessageID int64

func (id MessageID) I64() int64 { return int64(id) }

type Message struct {
	ID        MessageID
	ChatID    ChatID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
	WasRead   bool
}
