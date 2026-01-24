package user

import (
	"context"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type db interface {
	UserByLogin(ctx context.Context, login domain.Login) (*domain.User, error)
	UserByID(ctx context.Context, user domain.UserID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error)
	UpdateUser(ctx context.Context, updated *domain.User) error
	DeleteUser(ctx context.Context, user domain.UserID) error
}

type cache interface {
	WriteSession(ctx context.Context, session domain.Session, user domain.UserID, expiration time.Duration) error
	ClearSession(ctx context.Context, session domain.Session) error
	UserIDFromSession(ctx context.Context, session domain.Session, expiration time.Duration) domain.UserID
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db          db
	cache       cache
	transaction txManager

	expiration time.Duration
}

func New(dbAdapter db, cacheAdapter cache, transaction txManager) *Service {
	return &Service{
		db:          dbAdapter,
		cache:       cacheAdapter,
		transaction: transaction,
	}
}

func (s *Service) GetByCredentials(ctx context.Context, login domain.Login, rawPassword string) (*domain.User, error) {
	user, err := s.db.UserByLogin(ctx, login)
	if errors.Is(errors.NotFound, err) {
		return nil, errors.E(err).Kind(errors.InvalidRequest).Message("invalid login")
	}
	if err != nil {
		return nil, errors.E(err).Debug("s.db.UserByLogin")
	}

	if !user.Password.Equal(rawPassword, login) {
		return nil, errors.E(err).Kind(errors.InvalidRequest).Message("invalid password")
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error) {
	var resultID domain.UserID
	err := s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		id, err := s.db.CreateUser(ctx, user)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateUser")
		}
		resultID = id
		return nil
	})

	return resultID, err
}

func (s *Service) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.db.UpdateUser(ctx, user); err != nil {
			return errors.E(err).Debug("s.db.UpdateUser")
		}
		return nil
	})
}

func (s *Service) DeleteUser(ctx context.Context, id domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.db.DeleteUser(ctx, id); err != nil {
			return errors.E(err).Debug("s.db.DeleteUser")
		}
		return nil
	})
}

func (s *Service) Search(ctx context.Context, user domain.UserID, query string) ([]*domain.User, error) {
	// users, err := s.db.Users(ctx)
	// if err != nil {
	// 	return nil, errors.E(err).Debug("s.db.Chats")
	// }

	// query = strings.ToLower(query)

	// users = sliceutil.Filter(users, func(chat *domain.User) bool {
	// 	return
	// })

	// return users, nil
	return nil, nil
}
