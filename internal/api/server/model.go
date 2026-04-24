package server

import (
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
	"chattery/internal/utils/sliceutil"
)

type GetServersResponse struct {
	Servers []ServerResponse `json:"servers"`
}

type ServerResponse struct {
	Name   string          `json:"name"`
	Role   string          `json:"role"`
	Topics []TopicResponse `json:"topics"`
	ID     int64           `json:"id"`
}

type TopicResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ID   int64  `json:"id"`
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

type PostUpdateServerRequest struct {
	Name     string `json:"name"`
	ServerID int64  `json:"server_id"`
}

type PostCreateTopicRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ServerID int64  `json:"server_id"`
}

type PostCreateTopicResponse struct {
	ID int64 `json:"id"`
}

type PostUpdateTopicRequest struct {
	Name    string `json:"name"`
	TopicID int64  `json:"topic_id"`
}

type PostMessageRequest struct {
	Text    string `json:"text"`
	TopicID int64  `json:"topic_id"`
}

type GetTopicMessagesRequest struct {
	Cursor *Cursor `json:"cursor"`
}

type GetTopicMessagesResponse struct {
	Cursor   *Cursor   `json:"cursor"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Sender    UserInfo `json:"sender"`
	ID        int64    `json:"id"`
	SenderID  int64    `json:"sender_id"`
}

type UserInfo struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	ID       int64  `json:"id"`
}

type Cursor struct {
	Timestamp time.Time `json:"timestamp"`
	TopicID   int64     `json:"topic_id"`
	MessageID int64     `json:"message_id"`
}

type DeleteServerRequest struct {
	ServerID int64 `json:"server_id"`
}

type DeleteTopicRequest struct {
	TopicID int64 `json:"topic_id"`
}

func convertServerResponse(server *domain.Server) ServerResponse {
	return ServerResponse{
		ID:     server.ID.I64(),
		Name:   server.Name,
		Role:   server.Role.String(),
		Topics: sliceutil.Map(server.Topics, convertTopicResponse),
	}
}

func convertTopicResponse(topic *domain.Topic) TopicResponse {
	return TopicResponse{
		ID:   topic.ID.I64(),
		Name: topic.Name,
		Type: string(topic.Type),
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
		CreatedAt: render.Timestamp(msg.CreatedAt),
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
		Role:     domain.ServerRoleMember,
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
