package chat

import (
	"context"

	"chattery/internal/domain"
)

type db interface {
	AddParticipant(ctx context.Context, participant *domain.Participant) error
	DeleteParticipant(ctx context.Context, user domain.UserID, chat domain.ChatID) error
	ListParticipants(ctx context.Context, chat domain.ChatID) ([]*domain.Participant, error)

	Chats(ctx context.Context) ([]*domain.Chat, error)
	CreateChat(ctx context.Context, chat *domain.Chat) (domain.ChatID, error)
	DeleteChat(ctx context.Context, chat domain.ChatID) error

	UserChats(ctx context.Context, user domain.UserID) ([]*domain.Chat, error)

	CreateMessage(ctx context.Context, message *domain.Message) (domain.MessageID, error)
	FirstPageOfMessages(ctx context.Context, chat domain.ChatID) ([]*domain.Message, *domain.MessageCursor, error)
	NextPageOfMessages(ctx context.Context, chat domain.ChatID, cursor *domain.MessageCursor) ([]*domain.Message, *domain.MessageCursor, error)
}

type pubsub interface {
	SendMessage(ctx context.Context, chat domain.ChatID, message domain.Message) error
	Subscribe(ctx context.Context, chat domain.ChatID, dst chan<- *domain.Message)
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db          db
	transaction txManager
	pubsub      pubsub

	subs map[domain.UserID][]Subscription
}

func New(dbAdapter db, pubsubAdapter pubsub, transaction txManager) *Service {
	return &Service{
		db:          dbAdapter,
		pubsub:      pubsubAdapter,
		transaction: transaction,
	}
}

type callback func(ctx context.Context, event *domain.Event) error

type Subscription interface {
	GetUserID() domain.UserID
	GetSession() domain.Session
	Write(event domain.Event)
	SubscribeToEvent(ctx context.Context, type_ domain.EventType, callback callback)
}

func (s *Service) Register(sub Subscription) {
	s.subs[sub.GetUserID()] = nil
}
