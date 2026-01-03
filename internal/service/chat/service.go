package chat

import (
	"chattery/internal/domain"
	"context"
)

type db interface {
	AddParticipant(ctx context.Context, chatID domain.ChatID, username domain.Username) error
	Chats(ctx context.Context) ([]*domain.Chat, error)
	CreateChat(ctx context.Context, chat *domain.Chat) (domain.ChatID, error)
	CreateMessage(ctx context.Context, message *domain.Message) (domain.MessageID, error)
	FirstPageOfMessages(ctx context.Context, chatID domain.ChatID) ([]*domain.Message, *domain.ChatCursor, error)
	NextPagesOfMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.ChatCursor) ([]*domain.Message, *domain.ChatCursor, error)
}

type pubsub interface {
	SendMessage(ctx context.Context, chat domain.ChatID, message domain.Message) error
	Subscribe(ctx context.Context, chat domain.ChatID, dst chan<- *domain.Message)
}

type Service struct {
	db     db
	pubsub pubsub
}

func New(dbAdapter db, pubsubAdapter pubsub) *Service {
	return &Service{
		db:     dbAdapter,
		pubsub: pubsubAdapter,
	}
}
