package user

import (
	"context"
	"time"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

type db interface {
	UserByLogin(ctx context.Context, login domain.Email) (*domain.User, error)
	UserByID(ctx context.Context, user domain.UserID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error)
	UpdateUser(ctx context.Context, updated *domain.User) error
	DeleteUser(ctx context.Context, user domain.UserID) error
}

type cache interface {
	WriteSession(ctx context.Context, session domain.Session, user domain.UserID, expiration time.Duration) error
	ClearSession(ctx context.Context, session domain.Session) error
	UserIDFromSession(ctx context.Context, session domain.Session, expiration time.Duration, refreshBefore time.Duration) (domain.UserID, bool, error)
}

type txManager interface {
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type Service struct {
	db            db
	cache         cache
	transaction   txManager
	cookieDomain  string
	expiration    time.Duration
	refreshBefore time.Duration
	debug         bool
}

func New(dbAdapter db, cacheAdapter cache, transaction txManager, cfg *config.Config) *Service {
	return &Service{
		db:            dbAdapter,
		cache:         cacheAdapter,
		transaction:   transaction,
		debug:         cfg.App.Debug,
		expiration:    cfg.Session.Expiration,
		refreshBefore: cfg.Session.RefreshBefore,
		cookieDomain:  cfg.Session.CookieDomain,
	}
}

func (s *Service) GetByCredentials(ctx context.Context, login domain.Email, rawPassword string) (*domain.User, error) {
	user, err := s.db.UserByLogin(ctx, login)
	if errutil.Is(errutil.NotFound, err) {
		return nil, errutil.E(err).Kind(errutil.Permission).Messagef("user with login %q not found", login.String())
	}
	if err != nil {
		return nil, errutil.E(err).Debug("s.db.UserByLogin")
	}

	if !user.Password.Equal(rawPassword, login) {
		return nil, errutil.E(err).Kind(errutil.Permission).Message("invalid password")
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := s.db.UserByID(ctx, userID)
	if errutil.Is(errutil.NotFound, err) {
		return nil, errutil.E(err).Kind(errutil.NotFound).Messagef("user with id %d not found", userID)
	}
	if err != nil {
		return nil, errutil.E(err).Debug("s.db.UserByID")
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error) {
	var resultID domain.UserID
	err := s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		id, err := s.db.CreateUser(ctx, user)
		if err != nil {
			return errutil.E(err).Debug("s.db.CreateUser")
		}
		resultID = id
		return nil
	})

	return resultID, err
}

func (s *Service) UpdateUser(ctx context.Context, user *domain.User, currentPassword string) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		existing, err := s.db.UserByID(ctx, user.ID)
		if err != nil {
			return errutil.E(err).Debug("s.db.UserByID")
		}

		updated := mergeUserUpdate(existing, user)
		loginChanged := updated.Login != existing.Login
		passwordChanged := user.Password.Plain() != ""

		if loginChanged || passwordChanged {
			if currentPassword == "" || !existing.Password.Equal(currentPassword, existing.Login) {
				return errutil.E().
					Kind(errutil.Permission).
					Message("current password is invalid")
			}
		}

		if passwordChanged {
			updated.Password = domain.NewPassword(user.Password.Plain(), updated.Login)
		} else if loginChanged {
			updated.Password = domain.NewPassword(currentPassword, updated.Login)
		}

		if err := s.db.UpdateUser(ctx, updated); err != nil {
			return errutil.E(err).Debug("s.db.UpdateUser")
		}
		return nil
	})
}

func mergeUserUpdate(existing *domain.User, requested *domain.User) *domain.User {
	updated := *existing

	if requested.Username != "" {
		updated.Username = requested.Username
	}
	if requested.Login != "" {
		updated.Login = requested.Login
	}

	return &updated
}

func (s *Service) DeleteUser(ctx context.Context, id domain.UserID) error {
	return s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		if err := s.db.DeleteUser(ctx, id); err != nil {
			return errutil.E(err).Debug("s.db.DeleteUser")
		}
		return nil
	})
}
