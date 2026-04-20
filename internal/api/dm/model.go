package dm_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

type GetDMsResponse struct {
	DMs []DM `json:"dms"`
}

type DM struct {
	ID      int64            `json:"id"`
	Unread  bool             `json:"unread"`
	User    UserInfo         `json:"user"`
	Message *LastMessageInfo `json:"message,omitempty"`
}

type LastMessageInfo struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

type PostCreateDMRequest struct {
	ParticipantID int64 `json:"participant_id"`
}

type PostCreateDMResponse struct {
	ID int64 `json:"id"`
}

type PostMessageRequest struct {
	DMID int64  `json:"dm_id"`
	Text string `json:"text"`
}

type GetMessagesRequest struct {
	Cursor *Cursor `json:"cursor"`
}

type GetMessagesResponse struct {
	Messages []*Message `json:"messages"`
	Cursor   *Cursor    `json:"cursor"`
}

type Message struct {
	ID        int64    `json:"id"`
	Sender    UserInfo `json:"sender"`
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type Cursor struct {
	DMID      int64     `json:"dm_id"`
	MessageID int64     `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
}

func convertDMResponse(dm *domain.DM, users map[domain.UserID]*domain.User) DM {
	user := UserInfo{}
	if u, ok := users[dm.OtherParticipantID]; ok {
		user = UserInfo{
			ID:       u.ID.I64(),
			Username: u.Username.String(),
			Avatar:   "/v1/image/" + u.Username.String() + ".png",
		}
	}

	return DM{
		ID:      dm.ID.I64(),
		Unread:  dm.LastMessage.ID > dm.LastReadMessageID,
		User:    user,
		Message: convertLastMessageResponse(dm.LastMessage),
	}
}

func convertGetDMsResponse(dms []*domain.DM, users map[domain.UserID]*domain.User) *GetDMsResponse {
	dmResponses := make([]DM, 0, len(dms))
	for _, dm := range dms {
		dmResponses = append(dmResponses, convertDMResponse(dm, users))
	}
	return &GetDMsResponse{
		DMs: dmResponses,
	}
}

func convertMessageResponse(msg *domain.DMMessage, users map[domain.UserID]*domain.User) *Message {
	sender := UserInfo{}
	if user, ok := users[msg.SenderID]; ok {
		sender = UserInfo{
			ID:       user.ID.I64(),
			Username: user.Username.String(),
			// Avatar:   user.AvatarID.String(),
			Avatar: "/v1/image/" + user.Username.String() + ".png",
		}
	}
	return &Message{
		ID:        msg.ID.I64(),
		Sender:    sender,
		Text:      msg.Text,
		CreatedAt: render.Timestamp(msg.CreatedAt),
	}
}

func convertGetMessagesRequest(request *GetMessagesRequest) *domain.DMCursor {
	return &domain.DMCursor{
		ChatID:    domain.DMID(request.Cursor.DMID),
		MessageID: domain.DMMessageID(request.Cursor.MessageID),
		Timestamp: request.Cursor.Timestamp,
	}
}

func convertPostCreateDMResponse(id domain.DMID) *PostCreateDMResponse {
	return &PostCreateDMResponse{
		ID: id.I64(),
	}
}

func convertCursorResponse(cursor *domain.DMCursor) *Cursor {
	if cursor == nil {
		return nil
	}
	return &Cursor{
		DMID:      cursor.ChatID.I64(),
		MessageID: cursor.MessageID.I64(),
		Timestamp: cursor.Timestamp,
	}
}

func convertGetMessagesResponse(cursor *domain.DMCursor, messages []*domain.DMMessage, users map[domain.UserID]*domain.User) *GetMessagesResponse {
	msgs := make([]*Message, 0, len(messages))
	for _, msg := range messages {
		msgs = append(msgs, convertMessageResponse(msg, users))
	}
	return &GetMessagesResponse{
		Messages: msgs,
		Cursor:   convertCursorResponse(cursor),
	}
}

func convertPostMessageRequest(request *PostMessageRequest, userID domain.UserID) *domain.DMMessage {
	return &domain.DMMessage{
		DMID:     domain.DMID(request.DMID),
		SenderID: userID,
		Text:     request.Text,
	}
}

func convertLastMessageResponse(message domain.DMMessage) *LastMessageInfo {
	if message.ID == 0 {
		return nil
	}
	return &LastMessageInfo{
		Date:    render.Timestamp(message.CreatedAt),
		Content: message.Text,
	}
}
