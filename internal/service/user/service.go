package user

import (
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"context"
)

type userDB interface {
	UserByLogin(ctx context.Context, login domain.Login) (*domain.User, error)
}

type Service struct {
	db userDB
}

func New(db userDB) *Service {
	return &Service{
		db: db,
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
