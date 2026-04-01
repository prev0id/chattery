package server_adapter

import (
	"context"

	"chattery/internal/client/postgres"
	"chattery/internal/domain"
	"chattery/internal/utils/database"
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

func (a *Adapter) UpdateServer(ctx context.Context, server *domain.Server) error {
	req := &postgres.UpdateServerParams{
		ID:   server.ID.I64(),
		Name: server.Name,
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

func (a *Adapter) CreateServerParticipant(ctx context.Context, participant *domain.ServerParticipant) error {
	req := &postgres.CreateServerParticipantParams{
		ServerID: participant.ServerID.I64(),
		UserID:   participant.UserID.I64(),
		Role:     participant.Role.String(),
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

func (a *Adapter) CreateTopic(ctx context.Context, topic *domain.Topic) (domain.TopicID, error) {
	req := &postgres.CreateTopicParams{
		ServerID: topic.ServerID.I64(),
		Name:     topic.Name,
		Type:     topic.Type.String(),
	}

	id, err := a.db.Query(ctx).CreateTopic(ctx, req)
	if err != nil {
		return 0, errors.E(err).Debug("Query.CreateTopic")
	}
	return domain.TopicID(id), nil
}

func (a *Adapter) UpdateTopic(ctx context.Context, topic *domain.Topic) error {
	req := &postgres.UpdateTopicParams{
		ID:   topic.ID.I64(),
		Name: topic.Name,
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

func (a *Adapter) CreateMessage(ctx context.Context, message *domain.TopicMessage) error {
	req := &postgres.CreateMessageParams{
		TopicID: message.TopicID.I64(),
		UserID:  message.SenderID.I64(),
		Text:    message.Text,
	}

	if _, err := a.db.Query(ctx).CreateMessage(ctx, req); err != nil {
		return errors.E(err).Debug("Query.CreateMessage")
	}
	return nil
}

func (a *Adapter) FirstPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error) {
	req := &postgres.FirstPageOfTopicMessagesParams{
		TopicID: cursor.ChatID.I64(),
		Limit:   int32(cursor.Limit),
	}

	messages, err := a.db.Query(ctx).FirstPageOfTopicMessages(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.FirstPageOfTopicMessages")
	}

	return sliceutil.Map(messages, convertServerMessageFromDB), nil
}

func (a *Adapter) NextPageOfTopicMessages(ctx context.Context, cursor *domain.TopicCursor) ([]*domain.TopicMessage, error) {
	req := &postgres.NextPagesOfTopicMessagesParams{
		TopicID:   cursor.ChatID.I64(),
		CreatedAt: cursor.Timestamp,
		ID:        cursor.MessageID.I64(),
		Limit:     int32(cursor.Limit),
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

func (a *Adapter) GetServer(ctx context.Context, serverID domain.ServerID) (*domain.Server, error) {
	rows, err := a.db.Query(ctx).GetServer(ctx, serverID.I64())
	if database.NotFound(err) {
		return nil, errors.E(err).Kind(errors.NotFound)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetServer")
	}

	return convertServerWithTopicsFromDB(rows), nil
}

func (a *Adapter) GetServerParticipant(ctx context.Context, serverID domain.ServerID, userID domain.UserID) (*domain.ServerParticipant, error) {
	req := &postgres.GetServerParticipantParams{
		ServerID: serverID.I64(),
		UserID:   userID.I64(),
	}
	participant, err := a.db.Query(ctx).GetServerParticipant(ctx, req)
	if database.NotFound(err) {
		return nil, errors.E(err).Kind(errors.NotFound)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetServerParticipant")
	}

	return convertServerParticipantFromDB(participant), nil
}

func (a *Adapter) GetTopic(ctx context.Context, topicID domain.TopicID) (*domain.Topic, error) {
	topic, err := a.db.Query(ctx).GetTopic(ctx, topicID.I64())
	if database.NotFound(err) {
		return nil, errors.E(err).Kind(errors.NotFound)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetTopic")
	}

	return convertTopicFromDB(topic), nil
}
