package server_adapter

import (
	"context"
	"time"

	"chattery/internal/client/postgres"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/sliceutil"
)

type queryProvider interface {
	Query(ctx context.Context) postgres.Querier
}

type Adapter struct {
	db queryProvider
}

func New(db queryProvider) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) CreateServer(ctx context.Context, name string) (domain.ServerID, error) {
	id, err := a.db.Query(ctx).CreateServer(ctx, name)
	if err != nil {
		return domain.ServerID(0), errors.E(err).Debug("Query.CreateServer")
	}
	return domain.ServerID(id), nil
}

func (a *Adapter) UpdateServer(ctx context.Context, serverID domain.ServerID, name string) error {
	req := &postgres.UpdateServerParams{
		ID:   serverID.I64(),
		Name: name,
	}

	err := a.db.Query(ctx).UpdateServer(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.UpdateServer")
	}
	return nil
}

func (a *Adapter) DeleteServer(ctx context.Context, serverID domain.ServerID) error {
	err := a.db.Query(ctx).DeleteServerByServerID(ctx, serverID.I64())
	if err != nil {
		return errors.E(err).Debug("Query.DeleteServerByServerID")
	}
	return nil
}

func (a *Adapter) CreateServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID, role domain.ServerRole) error {
	req := &postgres.CreateServerParticipantParams{
		ServerID: serverID.I64(),
		UserID:   userID.I64(),
		Role:     role.String(),
	}

	err := a.db.Query(ctx).CreateServerParticipant(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.CreateServerParticipant")
	}
	return nil
}

func (a *Adapter) DeleteServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	req := &postgres.DeleteServerParticipantParams{
		ServerID: serverID.I64(),
		UserID:   userID.I64(),
	}

	err := a.db.Query(ctx).DeleteServerParticipant(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.DeleteServerParticipant")
	}
	return nil
}

func (a *Adapter) CreateTopic(ctx context.Context, serverID domain.ServerID, name string, topicType domain.TopicType) (*domain.Topic, error) {
	req := &postgres.CreateTopicParams{
		ServerID: serverID.I64(),
		Name:     name,
		Type:     string(topicType),
	}

	topic, err := a.db.Query(ctx).CreateTopic(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.CreateTopic")
	}
	return convertTopicFromDB(topic), nil
}

func (a *Adapter) UpdateTopic(ctx context.Context, topicID domain.TopicID, name string) error {
	req := &postgres.UpdateTopicParams{
		ID:   topicID.I64(),
		Name: name,
	}

	err := a.db.Query(ctx).UpdateTopic(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.UpdateTopic")
	}
	return nil
}

func (a *Adapter) DeleteTopic(ctx context.Context, topicID domain.TopicID) error {
	err := a.db.Query(ctx).DeleteTopic(ctx, topicID.I64())
	if err != nil {
		return errors.E(err).Debug("Query.DeleteTopic")
	}
	return nil
}

func (a *Adapter) CreateMessage(ctx context.Context, topicID domain.TopicID, userID domain.UserID, text string) (domain.TopicMessageID, error) {
	req := &postgres.CreateMessageParams{
		TopicID: topicID.I64(),
		UserID:  userID.I64(),
		Text:    text,
	}

	id, err := a.db.Query(ctx).CreateMessage(ctx, req)
	if err != nil {
		return domain.TopicMessageID(0), errors.E(err).Debug("Query.CreateMessage")
	}
	return domain.TopicMessageID(id), nil
}

func (a *Adapter) FirstPageOfTopicMessages(ctx context.Context, topicID domain.TopicID, limit int) ([]*domain.TopicMessage, error) {
	req := &postgres.FirstPageOfTopicMessagesParams{
		TopicID: topicID.I64(),
		Limit:   int32(limit),
	}

	messages, err := a.db.Query(ctx).FirstPageOfTopicMessages(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.FirstPageOfTopicMessages")
	}

	return sliceutil.Map(messages, convertServerMessageFromDB), nil
}

func (a *Adapter) NextPagesOfTopicMessages(ctx context.Context, topicID domain.TopicID, createdAt time.Time, lastMessageID domain.TopicMessageID, limit int) ([]*domain.TopicMessage, error) {
	req := &postgres.NextPagesOfTopicMessagesParams{
		TopicID:   topicID.I64(),
		CreatedAt: createdAt,
		ID:        lastMessageID.I64(),
		Limit:     int32(limit),
	}

	messages, err := a.db.Query(ctx).NextPagesOfTopicMessages(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.NextPagesOfTopicMessages")
	}

	return sliceutil.Map(messages, convertServerMessageFromDB), nil
}

func (a *Adapter) GetUserServers(ctx context.Context, userID domain.UserID) ([]*domain.Server, error) {
	servers, err := a.db.Query(ctx).GetUserServers(ctx, userID.I64())
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetUserServers")
	}

	return convertServersFromDB(servers), nil
}
