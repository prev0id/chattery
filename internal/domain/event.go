package domain

type Event struct {
	Type    EventType
	Message *Message
	Chat    ChatID
}

type EventType int8

const (
	EventTypeSendMessage EventType = iota
	EventTypeRecieveMessage
	EventTypeJoinChat
	EventTypeLeaveChat
)
