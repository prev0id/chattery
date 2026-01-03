package chatadapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/sliceutil"
	"context"
)

type queryProvider interface {
	Query(ctx context.Context) postgres.Querier
}

type Adapter struct {
	db    queryProvider
	limit int
}

func New(cfg *config.Config, db queryProvider) *Adapter {
	return &Adapter{
		db:    db,
		limit: cfg.Chat.MessagesLimit,
	}
}

func (a *Adapter) AddParticipant(ctx context.Context, chatID domain.ChatID, username domain.Username) error {
	req := &postgres.AddParticipantParams{
		ChatID:   chatID.I64(),
		Username: username.String(),
	}

	if err := a.db.Query(ctx).AddParticipant(ctx, req); err != nil {
		return errors.E(err).Debug("a.db.Query.AddParticipant")
	}

	return nil
}

func (a *Adapter) Chats(ctx context.Context) ([]*domain.Chat, error) {
	chats, err := a.db.Query(ctx).Chats(ctx)
	if err != nil {
		return nil, errors.E(err).Debug("Query.Chats")
	}

	return sliceutil.Map(chats, convertChat), nil
}

func (a *Adapter) CreateChat(ctx context.Context, chat *domain.Chat) (domain.ChatID, error) {
	id, err := a.db.Query(ctx).CreateChat(ctx, chat.Type.String())
	if err != nil {
		return 0, errors.E(err).Debug("Query.CreateChat")
	}

	return domain.ChatID(id), nil
}

func (a *Adapter) CreateMessage(ctx context.Context, message *domain.Message) (domain.MessageID, error) {
	req := &postgres.CreateMessageParams{
		ChatID:   message.ChatID.I64(),
		Username: message.Sender.String(),
		Text:     message.Text,
	}

	id, err := a.db.Query(ctx).CreateMessage(ctx, req)
	if err != nil {
		return 0, errors.E(err).Debug("Query.CreateChat")
	}

	return domain.MessageID(id), nil
}

func (a *Adapter) FirstPageOfMessages(ctx context.Context, chatID domain.ChatID) ([]*domain.Message, *domain.ChatCursor, error) {
	req := &postgres.FirstPageOfMessagesParams{
		ChatID: chatID.I64(),
		Limit:  int32(a.limit) + 1,
	}

	msgs, err := a.db.Query(ctx).FirstPageOfMessages(ctx, req)
	if err != nil {
		return nil, nil, errors.E(err).Debug("Query.FirstPageOfMessages")
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

	msgs, err := a.db.Query(ctx).NextPagesOfMessages(ctx, req)
	if err != nil {
		return nil, nil, errors.E(err).Debug("Query.NextPagesOfMessages")
	}

	if len(msgs) != a.limit+1 {
		return sliceutil.Map(msgs[:a.limit], convertMessage), convertCursor(msgs[a.limit]), nil
	}

	return sliceutil.Map(msgs, convertMessage), nil, nil
}
