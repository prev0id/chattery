package server_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type ListServersResponse struct {
	Servers []*ServerResponse `json:"servers"`
}

type ServerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateServerRequest struct {
	Name string `json:"name"`
}

type CreateServerResponse struct {
	ID int64 `json:"id"`
}

type JoinServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type LeaveServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type UpdateServerRequest struct {
	ServerID int64  `json:"server_id"`
	Name     string `json:"name"`
}

type CreateTopicRequest struct {
	ServerID int64            `json:"server_id"`
	Name     string           `json:"name"`
	Type     domain.TopicType `json:"type"`
}

type CreateTopicResponse struct {
	ID int64 `json:"id"`
}

type UpdateTopicRequest struct {
	TopicID int64  `json:"topic_id"`
	Name    string `json:"name"`
}

type CreateTopicMessageRequest struct {
	TopicID int64  `json:"topic_id"`
	Text    string `json:"text"`
}

type CreateTopicMessageResponse struct {
	ID int64 `json:"id"`
}

type ListTopicMessagesRequest struct {
	Cursor *TopicCursor `json:"cursor"`
}

type ListTopicMessagesResponse struct {
	Messages []*TopicMessage `json:"messages"`
	Cursor   *TopicCursor    `json:"cursor"`
}

type TopicMessage struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type TopicCursor struct {
	TopicID   int64     `json:"topic_id"`
	MessageID int64     `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
}

type DeleteServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type DeleteTopicRequest struct {
	TopicID int64 `json:"topic_id"`
}

func convertServerResponse(server *domain.Server) *ServerResponse {
	return &ServerResponse{
		ID:   server.ID.I64(),
		Name: server.Name,
	}
}

func convertListServersResponse(servers []*domain.Server) *ListServersResponse {
	return &ListServersResponse{
		Servers: sliceutil.Map(servers, convertServerResponse),
	}
}

func convertCreateServerResponse(id domain.ServerID) *CreateServerResponse {
	return &CreateServerResponse{
		ID: id.I64(),
	}
}

func convertTopicResponse(topic *domain.Topic) *CreateTopicResponse {
	return &CreateTopicResponse{
		ID: topic.ID.I64(),
	}
}

func convertTopicCursorRequest(request *TopicCursor) *domain.TopicCursor {
	if request == nil {
		return nil
	}
	return &domain.TopicCursor{
		ChatID:    domain.TopicID(request.TopicID),
		MessageID: domain.TopicMessageID(request.MessageID),
		Timestamp: request.Timestamp,
	}
}

func convertTopicMessageResponse(msg *domain.TopicMessage) *TopicMessage {
	return &TopicMessage{
		ID:        msg.ID.I64(),
		SenderID:  msg.SenderID.I64(),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func convertTopicCursorResponse(cursor *domain.TopicCursor) *TopicCursor {
	if cursor == nil {
		return nil
	}
	return &TopicCursor{
		TopicID:   cursor.ChatID.I64(),
		MessageID: cursor.MessageID.I64(),
		Timestamp: cursor.Timestamp,
	}
}

func convertListTopicMessagesResponse(cursor *domain.TopicCursor, messages []*domain.TopicMessage) *ListTopicMessagesResponse {
	return &ListTopicMessagesResponse{
		Messages: sliceutil.Map(messages, convertTopicMessageResponse),
		Cursor:   convertTopicCursorResponse(cursor),
	}
}

func convertCreateTopicMessageResponse(id domain.TopicMessageID) *CreateTopicMessageResponse {
	return &CreateTopicMessageResponse{
		ID: id.I64(),
	}
}
