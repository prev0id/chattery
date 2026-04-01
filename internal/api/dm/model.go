package dm_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type GetDMsResponse struct {
	DMs []DM `json:"dms"`
}

type DM struct {
	ID          int64    `json:"id"`
	LastMessage *Message `json:"last_message,omitempty"`
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
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Cursor struct {
	DMID      int64     `json:"dm_id"`
	MessageID int64     `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
}

func convertDMResponse(dm *domain.DM) DM {
	return DM{
		ID:          dm.ID.I64(),
		LastMessage: convertMessageResponse(&dm.LastMessage),
	}
}

func convertGetDMsResponse(dms []*domain.DM) *GetDMsResponse {
	return &GetDMsResponse{
		DMs: sliceutil.Map(dms, convertDMResponse),
	}
}

func convertMessageResponse(msg *domain.DMMessage) *Message {
	return &Message{
		ID:        msg.ID.I64(),
		SenderID:  msg.SenderID.I64(),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
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

func convertGetMessagesResponse(cursor *domain.DMCursor, messages []*domain.DMMessage) *GetMessagesResponse {
	return &GetMessagesResponse{
		Messages: sliceutil.Map(messages, convertMessageResponse),
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
