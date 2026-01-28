package chat_api

import (
	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type CreatePublicChatRequest struct {
	Name string `json:"id"`
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

type SearchChatsResponse struct {
	Chats []Chat `json:"chats"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func convertChat(chat *domain.Chat) Chat {
	return Chat{
		ID:   chat.ID.I64(),
		Name: chat.Name,
		Type: chat.Type.String(),
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
