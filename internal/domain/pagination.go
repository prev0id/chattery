package domain

import "time"

type Cursor[ID any] struct {
	ID        ID
	Timestamp time.Time
}

type (
	TopicMessageCursor = Cursor[TopicMessageID]
)
