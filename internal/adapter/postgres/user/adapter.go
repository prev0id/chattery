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

func (a *Adapter) UserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := a.db.Query(ctx).UserByUsername(ctx, userID.I64())
	if database.NotFound(err) {
		return nil, errors.E(err).
			Kind(errors.NotFound).
			Messagef("user %d not found", userID)
	}
	if err != nil {
		return nil, errors.E(err).Debug("Query.UserByUsername")
	}
	return convertUserFromDB(user), nil
}

func (a *Adapter) CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error) {
	req := &postgres.CreateUserParams{
		Login:    user.Login.String(),
		Password: user.Password,
		Username: user.Username.String(),
	}
	id, err := a.db.Query(ctx).CreateUser(ctx, req)
	if database.IsConstraintViolation(err, uniqueLoginConstraint) {
		return domain.UserIsUnknown, errors.E(err).
			Kind(errors.Exist).
			Messagef("user with login %s already exists", user.Login)
	}
	if database.IsConstraintViolation(err, uniqueUsernameConstraint) {
		return domain.UserIsUnknown, errors.E(err).
			Kind(errors.Exist).
			Messagef("user with username %s already exists", user.Username)
	}
	if err != nil {
		return domain.UserIsUnknown, errors.E(err).Debug("Query.CreateUser")
	}

	return domain.UserID(id), nil
}

func (a *Adapter) UpdateUser(ctx context.Context, updated *domain.User) error {
	req := &postgres.UpdateUserParams{
		ID:       updated.ID.I64(),
		Login:    updated.Login.String(),
		Password: updated.Password,
		Username: updated.Username.String(),
		AvatarID: updated.AvatarID.String(),
	}

	err := a.db.Query(ctx).UpdateUser(ctx, req)
	if database.IsConstraintViolation(err, uniqueLoginConstraint) {
		return errors.E(err).
			Kind(errors.Exist).
			Messagef("user with login %s already exists", updated.Login)
	}
	if database.IsConstraintViolation(err, uniqueUsernameConstraint) {
		return errors.E(err).
			Kind(errors.Exist).
			Messagef("user with username %s already exists", updated.Username)
	}
	if err != nil {
		return errors.E(err).Debug("Query.UpdateUser")
	}

	return nil
}

func (a *Adapter) DeleteUser(ctx context.Context, userID domain.UserID) error {
	rows, err := a.db.Query(ctx).DeleteUserByID(ctx, userID.I64())
	if err != nil {
		return errors.E(err).Debug("Query.DeleteUserByID")
	}
	if rows == 0 {
		return errors.E().Kind(errors.NotFound).Messagef("user %d not found", userID)
	}
	return nil
}
