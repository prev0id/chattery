package chatadapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
)

func convertChat(chat *postgres.Chat) *domain.Chat {
	return &domain.Chat{
		ID:   domain.ChatID(chat.ID),
		Type: domain.ChatType(chat.Type),
	}
}

func convertMessage(message *postgres.ChatMessage) *domain.Message {
	return &domain.Message{
		ID:     domain.MessageID(message.ID),
		ChatID: domain.ChatID(message.ChatID),
		Sender: domain.Username(message.Username),
		Text:   message.Text,
	}
}

func convertCursor(message *postgres.ChatMessage) *domain.ChatCursor {
	return &domain.ChatCursor{
		ID:        domain.MessageID(message.ID),
		Timestamp: message.CreatedAt,
	}
}
