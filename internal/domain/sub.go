package domain

import "context"

type Subscriber interface {
	GetSession() Session
	GetUserID() UserID
	WriteMessage(ctx context.Context, message *TopicMessage) error
}
