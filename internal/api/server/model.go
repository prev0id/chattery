package server_api

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type GetServersResponse struct {
	Servers []ServerResponse `json:"servers"`
}

type ServerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PostCreateServerRequest struct {
	Name string `json:"name"`
}

type PostCreateServerResponse struct {
	ID int64 `json:"id"`
}

type PostJoinServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type PostLeaveServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type PostServerUpdateRequest struct {
	ServerID int64  `json:"server_id"`
	Name     string `json:"name"`
}

type PostCreateTopicRequest struct {
	ServerID int64  `json:"server_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type PostCreateTopicResponse struct {
	ID int64 `json:"id"`
}

type PostUpdateTopicRequest struct {
	TopicID int64  `json:"topic_id"`
	Name    string `json:"name"`
}

type PostMessageRequest struct {
	TopicID int64  `json:"topic_id"`
	Text    string `json:"text"`
}

type GetTopicMessagesRequest struct {
	Cursor *Cursor `json:"cursor"`
}

type GetTopicMessagesResponse struct {
	Messages []Message `json:"messages"`
	Cursor   *Cursor   `json:"cursor"`
}

type Message struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Sender    UserInfo  `json:"sender"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type Cursor struct {
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

func convertServerResponse(server *domain.Server) ServerResponse {
	return ServerResponse{
		ID:   server.ID.I64(),
		Name: server.Name,
	}
}

func convertGetServersResponse(servers []*domain.Server) *GetServersResponse {
	return &GetServersResponse{
		Servers: sliceutil.Map(servers, convertServerResponse),
	}
}

func convertPostCreateServerResponse(id domain.ServerID) *PostCreateServerResponse {
	return &PostCreateServerResponse{
		ID: id.I64(),
	}
}

func convertPostCreateTopicResponse(topicID domain.TopicID) *PostCreateTopicResponse {
	return &PostCreateTopicResponse{
		ID: topicID.I64(),
	}
}

func convertTopicCursorRequest(request *Cursor) *domain.TopicCursor {
	if request == nil {
		return nil
	}
	return &domain.TopicCursor{
		ChatID:    domain.TopicID(request.TopicID),
		MessageID: domain.TopicMessageID(request.MessageID),
		Timestamp: request.Timestamp,
	}
}

func convertMessageResponse(msg *domain.TopicMessage, users map[domain.UserID]*domain.User) Message {
	sender := UserInfo{}
	if user, ok := users[msg.SenderID]; ok {
		sender = UserInfo{
			ID:       user.ID.I64(),
			Username: user.Username.String(),
			// Avatar:   user.AvatarID.String(),
			Avatar: "/v1/image/" + user.Username.String() + ".png",
		}
	}
	return Message{
		ID:        msg.ID.I64(),
		SenderID:  msg.SenderID.I64(),
		Sender:    sender,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func convertCursorResponse(cursor *domain.TopicCursor) *Cursor {
	if cursor == nil {
		return nil
	}
	return &Cursor{
		TopicID:   cursor.ChatID.I64(),
		MessageID: cursor.MessageID.I64(),
		Timestamp: cursor.Timestamp,
	}
}

func convertGetTopicMessagesResponse(cursor *domain.TopicCursor, messages []*domain.TopicMessage, users map[domain.UserID]*domain.User) *GetTopicMessagesResponse {
	msgs := make([]Message, len(messages))
	for i, msg := range messages {
		msgs[i] = convertMessageResponse(msg, users)
	}
	return &GetTopicMessagesResponse{
		Messages: msgs,
		Cursor:   convertCursorResponse(cursor),
	}
}

func convertPostMessageRequest(request *PostMessageRequest, userID domain.UserID) *domain.TopicMessage {
	return &domain.TopicMessage{
		TopicID:  domain.TopicID(request.TopicID),
		SenderID: userID,
		Text:     request.Text,
	}
}

func convertPostJoinServerRequest(request *PostJoinServerRequest, userID domain.UserID) *domain.ServerParticipant {
	return &domain.ServerParticipant{
		ServerID: domain.ServerID(request.ServerID),
		UserID:   userID,
		Role:     domain.ServerRoleUser,
	}
}

func convertPostCreateTopicRequest(request *PostCreateTopicRequest) *domain.Topic {
	return &domain.Topic{
		ServerID: domain.ServerID(request.ServerID),
		Name:     request.Name,
		Type:     domain.TopicType(request.Type),
	}
}
func convertPostUpdateTopicRequest(request *PostUpdateTopicRequest) *domain.Topic {
	return &domain.Topic{
		ID:   domain.TopicID(request.TopicID),
		Name: request.Name,
	}
}
