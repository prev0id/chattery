package chat_adapter

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
		Sender: domain.UserID(message.UserID),
		Text:   message.Text,
	}
}

func convertCursor(message *postgres.ChatMessage) *domain.MessageCursor {
	return &domain.MessageCursor{
		ID:        domain.MessageID(message.ID),
		Timestamp: message.CreatedAt,
	}
}

func convertParticipant(participant *postgres.ChatParticipant) *domain.Participant {
	return &domain.Participant{
		User: domain.UserID(participant.UserID),
		Chat: domain.ChatID(participant.ChatID),
		Role: domain.ChatRole(participant.Role),
	}
}
