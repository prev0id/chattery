package domain

import "time"

type Cursor[ID any] struct {
	ID        ID
	Timestamp time.Time
}

type (
	MessageCursor = Cursor[MessageID]
	ChatCursor    = Cursor[ChatID]
	UserCursor    = Cursor[UserID]
)
