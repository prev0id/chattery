package chat_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type CreatePublicChatRequest struct {
	Name string `json:"name"`
}

type CreatePrivateChatRequest struct {
	WithUserID int64 `json:"user_id"`
}

type CreateChatResponse struct {
	ID int64 `json:"id"`
}

type JoinRequest struct {
	ID int64 `json:"id"`
}

type LeaveRequest struct {
	ID int64 `json:"id"`
}

type MyChatsResponse struct {
	Private []Chat `json:"private"`
	Public  []Chat `json:"public"`
}

type PrivateChatsResponse struct {
	Chats []Chat `json:"chats"`
}

type PublicChatsResponse struct {
	Chats []Chat `json:"chats"`
}

type SearchChatsResponse struct {
	Chats []Chat `json:"chats"`
}

type Chat struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
}

type LastMessage struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func convertChat(chat *domain.Chat) Chat {
	return Chat{
		ID:          chat.ID.I64(),
		Name:        chat.Name,
		Type:        chat.Type.String(),
		LastMessage: nil,
	}
}

func convertLastMessage(message *domain.Message) *LastMessage {
	if message == nil {
		return nil
	}
	return &LastMessage{
		ID:        message.ID.I64(),
		UserID:    message.SenderID.I64(),
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	}
}

func convertChatPreview(preview *domain.ChatPreview) Chat {
	if preview == nil {
		return Chat{}
	}

	return Chat{
		ID:          preview.ID.I64(),
		Name:        preview.Name,
		Type:        preview.Type.String(),
		LastMessage: convertLastMessage(preview.LastMessage),
	}
}

func converMyChatsResponse(chats []*domain.Chat) MyChatsResponse {
	private := sliceutil.Filter(chats, func(chat *domain.Chat) bool {
		return chat.Type == domain.ChatTypePrivate
	})
	public := sliceutil.Filter(chats, func(chat *domain.Chat) bool {
		return chat.Type == domain.ChatTypePublic
	})

	return MyChatsResponse{
		Private: sliceutil.Map(private, convertChat),
		Public:  sliceutil.Map(public, convertChat),
	}
}
