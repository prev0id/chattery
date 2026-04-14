package dm_adapter

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

func (a *Adapter) UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error) {
	dms, err := a.db.Query(ctx).UserDMs(ctx, userID.I64())
	if err != nil {
		return nil, errors.E(err).Debug("Query.UserDMs")
	}

	return sliceutil.Map(dms, convertDMFromDB), nil
}

func (a *Adapter) CreateDM(ctx context.Context) (domain.DMID, error) {
	id, err := a.db.Query(ctx).CreateDM(ctx)
	if err != nil {
		return 0, errors.E(err).Debug("Query.CreateDM")
	}
	return domain.DMID(id), nil
}

func (a *Adapter) CreateDMParticipant(ctx context.Context, dmID domain.DMID, userID domain.UserID) error {
	req := &postgres.CreateDMParticipantParams{
		DmID:   dmID.I64(),
		UserID: userID.I64(),
	}

	err := a.db.Query(ctx).CreateDMParticipant(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.CreateDMParticipant")
	}
	return nil
}

func (a *Adapter) CreateDMMessage(ctx context.Context, message *domain.DMMessage) (domain.DMMessageID, error) {
	req := &postgres.CreateDMMessageParams{
		DmID:   message.DMID.I64(),
		UserID: message.SenderID.I64(),
		Text:   message.Text,
	}

	id, err := a.db.Query(ctx).CreateDMMessage(ctx, req)
	if err != nil {
		return domain.DMMessageID(0), errors.E(err).Debug("Query.CreateDMMessage")
	}
	return domain.DMMessageID(id), nil
}

func (a *Adapter) SetDMLastReadMessage(ctx context.Context, dmID domain.DMID, userID domain.UserID, messageID domain.DMMessageID) error {
	req := &postgres.SetDMLastReadMessageParams{
		DmID:              dmID.I64(),
		UserID:            userID.I64(),
		LastReadMessageID: messageID.I64(),
	}

	err := a.db.Query(ctx).SetDMLastReadMessage(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.SetDMLastReadMessage")
	}
	return nil
}

func (a *Adapter) SetLastMessageInDM(ctx context.Context, dmID domain.DMID, messageID domain.DMMessageID) error {
	req := &postgres.SetLastMessageInDMParams{
		ID:            dmID.I64(),
		LastMessageID: messageID.I64(),
	}

	err := a.db.Query(ctx).SetLastMessageInDM(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.SetLastMessageInDM")
	}
	return nil
}

func (a *Adapter) FirstPageOfDMMessages(ctx context.Context, cursor *domain.DMCursor) ([]*domain.DMMessage, error) {
	req := &postgres.FirstPageOfDMMessagesParams{
		DmID:  cursor.ChatID.I64(),
		Limit: int32(cursor.Limit),
	}

	messages, err := a.db.Query(ctx).FirstPageOfDMMessages(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.FirstPageOfDMMessages")
	}

	return sliceutil.Map(messages, convertDMMessageFromDB), nil
}

func (a *Adapter) NextPagesOfDMMessages(ctx context.Context, cursor *domain.DMCursor) ([]*domain.DMMessage, error) {
	req := &postgres.NextPagesOfDMMessagesParams{
		DmID:      cursor.ChatID.I64(),
		CreatedAt: cursor.Timestamp,
		ID:        cursor.MessageID.I64(),
		Limit:     int32(cursor.Limit),
	}

	messages, err := a.db.Query(ctx).NextPagesOfDMMessages(ctx, req)
	if err != nil {
		return nil, errors.E(err).Debug("Query.NextPagesOfDMMessages")
	}

	return sliceutil.Map(messages, convertDMMessageFromDB), nil
}

func (a *Adapter) GetDMParticipant(ctx context.Context, dmID domain.DMID, userID domain.UserID) (*domain.DMParticipant, error) {
	req := &postgres.GetDMParticipantParams{
		DmID:   dmID.I64(),
		UserID: userID.I64(),
	}

	participant, err := a.db.Query(ctx).GetDMParticipant(ctx, req)
	if database.NotFound(err) {
		return nil, errors.E(err).Kind(errors.NotFound)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetDMParticipant")
	}

	return convertDMParticipantFromDB(participant), nil
}

func (a *Adapter) GetDMBetweenUsers(ctx context.Context, userID1, userID2 domain.UserID) (*domain.DM, error) {
	req := &postgres.GetDMBetweenUsersParams{
		UserID:   userID1.I64(),
		UserID_2: userID2.I64(),
	}

	dmID, err := a.db.Query(ctx).GetDMBetweenUsers(ctx, req)
	if database.NotFound(err) {
		return nil, errors.E(err).Kind(errors.NotFound)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetDMBetweenUsers")
	}

	return &domain.DM{
		ID: domain.DMID(dmID),
	}, nil
}

func (a *Adapter) GetDMParticipants(ctx context.Context, dmID domain.DMID) ([]domain.UserID, error) {
	ids, err := a.db.Query(ctx).GetDMParticipants(ctx, dmID.I64())
	if err != nil {
		return nil, errors.E(err).Debug("Query.GetDMParticipants")
	}
	return sliceutil.Map(ids, func(id int64) domain.UserID { return domain.UserID(id) }), nil
}
