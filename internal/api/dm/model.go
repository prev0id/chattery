package dm_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type ListDMResponse struct {
	DMs []DM `json:"dms"`
}

type DM struct {
	ID          int64    `json:"id"`
	LastMessage *Message `json:"last_message,omitempty"`
}

type CreateDMRequest struct {
	ParticipantID int64 `json:"participant_id"`
}

type CreateDMResponse struct {
	ID int64 `json:"id"`
}

type CreateMessageRequest struct {
	DMID int64  `json:"dm_id"`
	Text string `json:"text"`
}

type ListMessagesRequest struct {
	Cursor *Cursor `json:"cursor"`
}

type ListMessagesResponse struct {
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
		LastMessage: convertDMMessageResponse(&dm.LastMessage),
	}
}

func convertListDMResponse(dms []*domain.DM) *ListDMResponse {
	return &ListDMResponse{
		DMs: sliceutil.Map(dms, convertDMResponse),
	}
}

func convertDMMessageResponse(msg *domain.DMMessage) *Message {
	return &Message{
		ID:        msg.ID.I64(),
		SenderID:  msg.SenderID.I64(),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func convertListMessagesRequest(request *ListMessagesRequest) *domain.DMCursor {
	return &domain.DMCursor{
		ChatID:    domain.DMID(request.Cursor.DMID),
		MessageID: domain.DMMessageID(request.Cursor.MessageID),
		Timestamp: request.Cursor.Timestamp,
	}
}

func convertCreateDMResponse(id domain.DMID) *CreateDMResponse {
	return &CreateDMResponse{
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

func convertListMessagesResponse(cursor *domain.DMCursor, messages []*domain.DMMessage) *ListMessagesResponse {
	return &ListMessagesResponse{
		Messages: sliceutil.Map(messages, convertDMMessageResponse),
		Cursor:   convertCursorResponse(cursor),
	}
}

func convertCreateMessageRequest(request *CreateMessageRequest, userID domain.UserID) *domain.DMMessage {
	return &domain.DMMessage{
		DMID:     domain.DMID(request.DMID),
		SenderID: userID,
		Text:     request.Text,
	}
}
