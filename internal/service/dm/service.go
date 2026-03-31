package dm

import (
	"context"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type db interface {
	UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error)
	CreateDM(ctx context.Context) (domain.DMID, error)
	CreateDMParticipant(ctx context.Context, dmID domain.DMID, userID domain.UserID) error
	CreateDMMessage(ctx context.Context, message *domain.DMMessage) (domain.DMMessageID, error)
	SetDMLastReadMessage(ctx context.Context, dmID domain.DMID, userID domain.UserID, messageID domain.DMMessageID) error
	SetLastMessageInDM(ctx context.Context, dmID domain.DMID, messageID domain.DMMessageID) error
	FirstPageOfDMMessages(ctx context.Context, cursor *domain.DMCursor) ([]*domain.DMMessage, error)
	NextPagesOfDMMessages(ctx context.Context, cursor *domain.DMCursor) ([]*domain.DMMessage, error)
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db          db
	transaction txManager
	limit       int
}

func New(dbAdapter db, transaction txManager, cfg *config.Config) *Service {
	return &Service{
		db:          dbAdapter,
		transaction: transaction,
		limit:       cfg.Chat.MessagesLimit,
	}
}

func (s *Service) CreateDM(ctx context.Context, participant1, participant2 domain.UserID) (domain.DMID, error) {
	var (
		dmID domain.DMID
		err  error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		dmID, err = s.db.CreateDM(ctx)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateDM")
		}

		if err = s.db.CreateDMParticipant(ctx, dmID, participant1); err != nil {
			return errors.E(err).Debug("s.db.CreateDMParticipant")
		}

		if err = s.db.CreateDMParticipant(ctx, dmID, participant2); err != nil {
			return errors.E(err).Debug("s.db.CreateDMParticipant")
		}

		return nil
	})

	return dmID, err
}

func (s *Service) UserDMs(ctx context.Context, userID domain.UserID) ([]*domain.DM, error) {
	dms, err := s.db.UserDMs(ctx, userID)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.UserDMs")
	}
	return dms, nil
}

func (s *Service) CreateDMMessage(ctx context.Context, message *domain.DMMessage) error {
	var (
		messageID domain.DMMessageID
		err       error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		messageID, err = s.db.CreateDMMessage(ctx, message)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateDMMessage")
		}

		if err = s.db.SetLastMessageInDM(ctx, message.DMID, messageID); err != nil {
			return errors.E(err).Debug("s.db.SetLastMessageInDM")
		}

		return nil
	})

	return err
}

func (s *Service) FirstPageOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	cursor.Limit = s.limit

	messages, err := s.db.FirstPageOfDMMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.FirstPageOfDMMessages")
	}

	if len(messages) > 0 {
		lastSeenMessage := messages[0]
		if err := s.db.SetDMLastReadMessage(ctx, cursor.ChatID, userID, lastSeenMessage.ID); err != nil {
			return nil, nil, errors.E(err).Debug("s.db.SetDMLastReadMessage")
		}
	}

	var nextCursor *domain.DMCursor

	if len(messages) == len(messages)-1 {
		lastMessage := messages[len(messages)-1]
		nextCursor = &domain.DMCursor{
			ChatID:    cursor.ChatID,
			MessageID: lastMessage.ID,
			Timestamp: lastMessage.CreatedAt,
		}
	}

	return messages, nextCursor, nil
}

func (s *Service) NextPagesOfDMMessages(ctx context.Context, userID domain.UserID, cursor *domain.DMCursor) ([]*domain.DMMessage, *domain.DMCursor, error) {
	cursor.Limit = s.limit

	messages, err := s.db.NextPagesOfDMMessages(ctx, cursor)
	if err != nil {
		return nil, nil, errors.E(err).Debug("s.db.NextPagesOfDMMessages")
	}

	var nextCursor *domain.DMCursor

	if len(messages) == len(messages)-1 {
		lastMessage := messages[len(messages)-1]
		nextCursor = &domain.DMCursor{
			ChatID:    cursor.ChatID,
			MessageID: lastMessage.ID,
			Timestamp: lastMessage.CreatedAt,
		}
	}

	return messages, nextCursor, nil
}
