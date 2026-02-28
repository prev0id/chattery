package chat_adapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
)

func convertChat(chat *postgres.Chat) *domain.Chat {
	return &domain.Chat{
		ID:   domain.ChatID(chat.ID),
		Name: chat.Name,
		Type: domain.ChatType(chat.Type),
	}
}

func convertMessage(message *postgres.ChatMessage) *domain.Message {
	return &domain.Message{
		ID:        domain.MessageID(message.ID),
		ChatID:    domain.ChatID(message.ChatID),
		SenderID:  domain.UserID(message.UserID),
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
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
		UserID: domain.UserID(participant.UserID),
		Chat:   domain.ChatID(participant.ChatID),
		Role:   domain.ChatRole(participant.Role),
	}
}

func convertChatPreview(row *postgres.UserChatPreviewByTypeRow) *domain.ChatPreview {
	preview := &domain.ChatPreview{
		ID:   domain.ChatID(row.ID),
		Name: row.Name,
		Type: domain.ChatType(row.Type),
	}

	if row.HasLastMessage {
		preview.LastMessage = &domain.Message{
			ID:        domain.MessageID(row.LastMessageID),
			ChatID:    domain.ChatID(row.ID),
			SenderID:  domain.UserID(row.LastMessageUserID),
			Text:      row.LastMessageText,
			CreatedAt: row.LastMessageCreatedAt,
		}
	}

	return preview
}
