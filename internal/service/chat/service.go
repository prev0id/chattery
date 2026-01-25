package chat

import (
	"context"
	"sync"

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

	subsMutex         sync.RWMutex
	chatSubsBySession map[domain.Session]context.CancelFunc
	subsByUserID      map[domain.UserID][]domain.Subscriber
}

func New(dbAdapter db, pubsubAdapter pubsub, transaction txManager) *Service {
	return &Service{
		db:                dbAdapter,
		pubsub:            pubsubAdapter,
		transaction:       transaction,
		chatSubsBySession: make(map[domain.Session]context.CancelFunc),
		subsByUserID:      make(map[domain.UserID][]domain.Subscriber),
	}
}
