package domain

import "time"

type ChatID int64

const ChatUnknown ChatID = 0

func (id ChatID) I64() int64 { return int64(id) }

type ChatType string

func (t ChatType) String() string { return string(t) }

const (
	ChatTypePersonal ChatType = "personal"
	ChatTypeGroup    ChatType = "group"
)

type Chat struct {
	ID           ChatID
	Type         ChatType
	Participants []Username
}

type MessageID int64

func (id MessageID) I64() int64 { return int64(id) }

type Message struct {
	ID        MessageID
	ChatID    ChatID
	Sender    Username
	Text      string
	CreatedAt time.Time
}

type ChatCursor struct {
	ID        MessageID
	Timestamp time.Time
}
