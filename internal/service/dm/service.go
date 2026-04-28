package dm

import (
	"context"

	"chattery/internal/config"
	"chattery/internal/domain"
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
	GetDMParticipant(ctx context.Context, dmID domain.DMID, userID domain.UserID) (*domain.DMParticipant, error)
	GetDMBetweenUsers(ctx context.Context, userID1, userID2 domain.UserID) (*domain.DM, error)
	GetDMParticipants(ctx context.Context, dmID domain.DMID) ([]domain.UserID, error)
	GetUserDMParticipantIDs(ctx context.Context, userID domain.UserID) ([]domain.UserID, error)
}

type redis interface {
	PublishToUser(ctx context.Context, userID domain.UserID, event []byte) error
}

type userCache interface {
	GetByID(userID domain.UserID) (*domain.User, error)
	List() []*domain.User
}

type dmCache interface {
	GetByID(dmID domain.DMID) (*domain.DM, error)
	List() []*domain.DM
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db          db
	transaction txManager
	redis       redis
	user        userCache
	cache       dmCache
	limit       int
}

func New(dbAdapter db, transaction txManager, redisAdapter redis, user userCache, cache dmCache, cfg *config.Config) *Service {
	return &Service{
		db:          dbAdapter,
		transaction: transaction,
		redis:       redisAdapter,
		limit:       cfg.Chat.MessagesLimit,
		user:        user,
		cache:       cache,
	}
}

func (s *Service) getNextCursor(dmID domain.DMID, messages []*domain.DMMessage) *domain.DMCursor {
	if len(messages) != s.limit {
		return nil
	}

	lastMessage := messages[len(messages)-1]
	return &domain.DMCursor{
		ChatID:    dmID,
		MessageID: lastMessage.ID,
		Timestamp: lastMessage.CreatedAt,
	}
}
