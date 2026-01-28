package subscriber

import (
	"time"

	"chattery/internal/domain"
)

type Event struct {
	Type    EventType `json:"type"`
	Message Message   `json:"message"`
	ChatID  int64     `json:"chat_id"`
}

type Message struct {
	Text      string    `json:"text"`
	SenderID  int64     `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}

type EventType string

const (
	EventTypeMessage   EventType = "message"
	EventTypeJoinChat  EventType = "join_chat"
	EventTypeLeaveChat EventType = "leave_chat"
)

func convertMessageToEvent(message *domain.Message) *Event {
	return &Event{
		Type:   EventTypeMessage,
		ChatID: message.ChatID.I64(),
		Message: Message{
			Text:      message.Text,
			SenderID:  message.SenderID.I64(),
			CreatedAt: message.CreatedAt,
		},
	}
}
func convertEventToMessage(event *Event, user domain.UserID) *domain.Message {
	return &domain.Message{
		ChatID:    domain.ChatID(event.ChatID),
		Text:      event.Message.Text,
		SenderID:  user,
		CreatedAt: time.Now(),
	}
}
