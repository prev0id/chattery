package user

import (
	"context"
	"time"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

type db interface {
	UserByLogin(ctx context.Context, login domain.Login) (*domain.User, error)
	UserByUsername(ctx context.Context, username domain.Username) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, username domain.Username, updated *domain.User) error
	DeleteUser(ctx context.Context, username domain.Username) error
}

type cache interface {
	WriteSession(ctx context.Context, session domain.Session, user domain.Username, expiration time.Duration) error
	ClearSession(ctx context.Context, session domain.Session) error
	UsernameFromSession(ctx context.Context, session domain.Session, expiration time.Duration) domain.Username
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

func (s *Service) ValidateCredentials(ctx context.Context, login domain.Login, rawPassword string) (*domain.User, error) {
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

func (s *Service) CreateUser(ctx context.Context, user *domain.User) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.validateLoginExists(ctx, user); err != nil {
			return errors.E(err).Debug("s.validateLoginExists")
		}

		if err := s.validateUsernameExists(ctx, user); err != nil {
			return errors.E(err).Debug("s.validateUsernameExists")
		}

		if err := s.db.CreateUser(ctx, user); err != nil {
			return errors.E(err).Debug("s.db.CreateUser")
		}
		return nil
	})
}

func (s *Service) UpdateUser(ctx context.Context, username domain.Username, updated *domain.User) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		stored, err := s.db.UserByUsername(ctx, username)
		if err != nil {
			return errors.E(err).Debug("s.db.UserByUsername")
		}

		if stored.Username != updated.Username {
			if err := s.validateUsernameExists(ctx, updated); err != nil {
				return errors.E(err).Debug("s.validateUsernameExists")
			}
		}

		if stored.Login != updated.Login {
			if err := s.validateLoginExists(ctx, updated); err != nil {
				return errors.E(err).Debug("s.validateLoginExists")
			}
		}

		if err := s.db.UpdateUser(ctx, username, updated); err != nil {
			return errors.E(err).Debug("s.db.UpdateUser")
		}
		return nil
	})
}

func (s *Service) DeleteUser(ctx context.Context, username domain.Username) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.db.DeleteUser(ctx, username); err != nil {
			return errors.E(err).Debug("s.db.DeleteUser")
		}
		return nil
	})
}

func (s *Service) validateLoginExists(ctx context.Context, user *domain.User) error {
	_, err := s.db.UserByLogin(ctx, user.Login)
	if !errors.Is(errors.NotFound, err) {
		return errors.E(err).Debug("s.db.UserByLogin")
	}
	if err == nil {
		return errors.E().Kind(errors.Exist).Messagef("user with login %s already exists", user.Login)
	}
	return nil
}

func (s *Service) validateUsernameExists(ctx context.Context, user *domain.User) error {
	_, err := s.db.UserByUsername(ctx, user.Username)
	if !errors.Is(errors.NotFound, err) {
		return errors.E(err).Debug("s.db.UserByUsername")
	}
	if err == nil {
		return errors.E().Kind(errors.Exist).Messagef("user with username %s already exists", user.Username)
	}
	return nil
}
