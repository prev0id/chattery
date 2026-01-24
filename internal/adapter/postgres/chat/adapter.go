package chat_adapter

import (
	"context"

	"chattery/internal/client/postgres"
	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/sliceutil"
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

func (a *Adapter) AddParticipant(ctx context.Context, participant *domain.Participant) error {
	req := &postgres.AddParticipantParams{
		ChatID: participant.Chat.I64(),
		UserID: participant.User.I64(),
		Role:   participant.Role.String(),
	}

	if err := a.db.Query(ctx).AddParticipant(ctx, req); err != nil {
		return errors.E(err).Debug("a.db.Query.AddParticipant")
	}
	return nil
}

func (a *Adapter) DeleteParticipant(ctx context.Context, user domain.UserID, chat domain.ChatID) error {
	req := &postgres.DeleteParticipantParams{
		ChatID: chat.I64(),
		UserID: user.I64(),
	}

	if err := a.db.Query(ctx).DeleteParticipant(ctx, req); err != nil {
		return errors.E(err).Debug("a.db.Query.DeleteParticipant")
	}
	return nil
}

func (a *Adapter) ListParticipants(ctx context.Context, chat domain.ChatID) ([]*domain.Participant, error) {
	participants, err := a.db.Query(ctx).ParticipantsForChat(ctx, chat.I64())
	if err != nil {
		return nil, errors.E(err).Debug("a.db.Query.ParticipantsForChat")
	}
	return sliceutil.Map(participants, convertParticipant), nil
}

func (a *Adapter) Chats(ctx context.Context) ([]*domain.Chat, error) {
	chats, err := a.db.Query(ctx).Chats(ctx)
	if err != nil {
		return nil, errors.E(err).Debug("Query.Chats")
	}

	return sliceutil.Map(chats, convertChat), nil
}

func (a *Adapter) CreateChat(ctx context.Context, chat *domain.Chat) (domain.ChatID, error) {
	req := &postgres.CreateChatParams{
		Type: chat.Type.String(),
		Name: chat.Name,
	}
	id, err := a.db.Query(ctx).CreateChat(ctx, req)
	if err != nil {
		return 0, errors.E(err).Debug("a.db.Query.CreateChat")
	}
	return domain.ChatID(id), nil
}

func (a *Adapter) DeleteChat(ctx context.Context, chat domain.ChatID) error {
	if err := a.db.Query(ctx).DeleteChat(ctx, chat.I64()); err != nil {
		return errors.E(err).Debug("a.db.Query.DeleteChat")
	}
	return nil
}

func (a *Adapter) UserChats(ctx context.Context, user domain.UserID) ([]*domain.Chat, error) {
	chats, err := a.db.Query(ctx).UserChats(ctx, user.I64())
	if err != nil {
		return nil, errors.E(err).Debug("a.db.Query.UserChats")
	}
	return sliceutil.Map(chats, convertChat), nil
}

func (a *Adapter) CreateMessage(ctx context.Context, message *domain.Message) (domain.MessageID, error) {
	req := &postgres.CreateMessageParams{
		ChatID: message.ChatID.I64(),
		UserID: message.Sender.I64(),
		Text:   message.Text,
	}

	id, err := a.db.Query(ctx).CreateMessage(ctx, req)
	if err != nil {
		return 0, errors.E(err).Debug("Query.CreateChat")
	}

	return domain.MessageID(id), nil
}

func (a *Adapter) FirstPageOfMessages(ctx context.Context, chatID domain.ChatID) ([]*domain.Message, *domain.MessageCursor, error) {
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

func (a *Adapter) NextPageOfMessages(ctx context.Context, chatID domain.ChatID, cursor *domain.MessageCursor) ([]*domain.Message, *domain.MessageCursor, error) {
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
