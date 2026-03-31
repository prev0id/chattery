package domain

import "time"

type Cursor[ChatID, MessageID any] struct {
	ChatID    ChatID
	MessageID MessageID
	Timestamp time.Time
	Limit     int
}

type (
	TopicCursor = Cursor[TopicID, TopicMessageID]
	DMCursor    = Cursor[DMID, DMMessageID]
)
