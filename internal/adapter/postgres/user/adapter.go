package useradapter

import (
	"context"

	"chattery/internal/client/postgres"
	"chattery/internal/domain"
	"chattery/internal/utils/database"
	"chattery/internal/utils/errors"
)

const (
	uniqueLoginConstraint    = "users_login_key"
	uniqueUsernameConstraint = "users_username_key"
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

func (a *Adapter) UserByLogin(ctx context.Context, login domain.Login) (*domain.User, error) {
	user, err := a.db.Query(ctx).UserByLogin(ctx, login.String())
	if database.NotFound(err) {
		return nil, errors.E(err).
			Kind(errors.NotFound).
			Messagef("user with login %s not found", login)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.UserByLogin")
	}
	return convertUserFromDB(user), nil
}
func (a *Adapter) UserByUsername(ctx context.Context, username domain.Username) (*domain.User, error) {
	user, err := a.db.Query(ctx).UserByUsername(ctx, username.String())
	if database.NotFound(err) {
		return nil, errors.E(err).
			Kind(errors.NotFound).
			Messagef("user %s not found", username)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.UserByUsername")
	}
	return convertUserFromDB(user), nil
}

func (a *Adapter) CreateUser(ctx context.Context, user *domain.User) error {
	req := &postgres.CreateUserParams{
		Login:    user.Login.String(),
		Password: user.Password,
		Username: user.Username.String(),
	}
	err := a.db.Query(ctx).CreateUser(ctx, req)
	if err != nil {
		return errors.E(err).Debug("Query.CreateUser")
	}
	return nil
}

func (a *Adapter) UpdateUser(ctx context.Context, username domain.Username, updated *domain.User) error {
	req := &postgres.UpdateUserParams{
		OldUsername: username.String(),
		NewLogin:    updated.Login.String(),
		NewPassword: updated.Password,
		NewUsername: updated.Username.String(),
		NewAvatarID: updated.AvatarID.String(),
	}
	if err := a.db.Query(ctx).UpdateUser(ctx, req); err != nil {
		return errors.E(err).Debug("Query.UpdateUser")
	}
	return nil
}
func (a *Adapter) DeleteUser(ctx context.Context, username domain.Username) error {
	rows, err := a.db.Query(ctx).DeleteUserByUsername(ctx, username.String())
	if err != nil {
		return errors.E(err).Debug("Query.DeleteUserByUsername")
	}
	if rows == 0 {
		return errors.E().Kind(errors.NotFound).Messagef("user %s not found", username)
	}
	return nil
}
