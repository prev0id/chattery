package chatadapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
	"context"
	"fmt"
)

type Adapter struct {
	db    postgres.Querier
	limit int
}

func New(cfg *config.Config, db postgres.Querier) *Adapter {
	return &Adapter{
		limit: cfg.MessagesLimit,
	}
}

func (a *Adapter) AddParticipant(ctx context.Context, chatID domain.ChatID, username domain.Username) error {
	req := &postgres.AddParticipantParams{
		ChatID:   chatID.I64(),
		Username: username.String(),
	}

	if err := a.db.AddParticipant(ctx, req); err != nil {
		return fmt.Errorf("a.db.AddParticipant: %w", err)
	}

	return nil
}

func (a *Adapter) Chats(ctx context.Context) ([]*domain.Chat, error) {
	chats, err := a.db.Chats(ctx)
	if err != nil {
		return nil, fmt.Errorf("a.db.Chats: %w", err)
	}

	return sliceutil.Map(chats, convertChat), nil
}

func (a *Adapter) CreateChat(ctx context.Context, chat *domain.Chat) (domain.ChatID, error) {
	id, err := a.db.CreateChat(ctx, chat.Type.String())
	if err != nil {
		return 0, fmt.Errorf("a.db.CreateChat: %w", err)
	}

	return domain.ChatID(id), nil
}

func (a *Adapter) CreateMessage(ctx context.Context, message *domain.Message) (domain.MessageID, error) {
	req := &postgres.CreateMessageParams{
		ChatID:   message.ChatID.I64(),
		Username: message.Sender.String(),
		Text:     message.Text,
	}

	id, err := a.db.CreateMessage(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("a.db.CreateChat: %w", err)
	}

	return domain.MessageID(id), nil
}

func (a *Adapter) FirstPageOfMessages(ctx context.Context, chatID domain.ChatID) ([]*domain.Message, *domain.ChatCursor, error) {
	req := &postgres.FirstPageOfMessagesParams{
		ChatID: chatID.I64(),
		Limit:  int32(a.limit) + 1,
	}

	msgs, err := a.db.FirstPageOfMessages(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("a.db.FirstPageOfMessages: %w", err)
	}

	if len(msgs) != a.limit+1 {
		return sliceutil.Map(msgs[:a.limit], convertMessage), convertCursor(msgs[a.limit]), nil
	}

	return sliceutil.Map(msgs, convertMessage), nil, nil
}

func (a *Adapter) NextPagesOfMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.ChatCursor) ([]*domain.Message, *domain.ChatCursor, error) {
	req := &postgres.NextPagesOfMessagesParams{
		ChatID:    chatID.I64(),
		ID:        cursor.ID.I64(),
		CreatedAt: cursor.Timestamp,
		Limit:     int32(a.limit) + 1,
	}

	msgs, err := a.db.NextPagesOfMessages(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("a.db.FirstPageOfMessages: %w", err)
	}

	if len(msgs) != a.limit+1 {
		return sliceutil.Map(msgs[:a.limit], convertMessage), convertCursor(msgs[a.limit]), nil
	}

	return sliceutil.Map(msgs, convertMessage), nil, nil
}
