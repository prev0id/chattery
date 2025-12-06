package postgres_adapter

import (
	"context"
	"errors"
	"fmt"

	postgres_client "chattery/backend/user_service/internal/client/postgres"
	"chattery/backend/user_service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueLoginConstraint    = "users_login_key"
	uniqueUsernameConstraint = "users_username_key"
)

type Adapter struct {
	db *postgres_client.Queries
}

func New(db *postgres_client.Queries) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) CreateUser(ctx context.Context, user domain.User) (domain.UserID, error) {
	id, err := a.db.InsertUser(ctx, postgres_client.InsertUserParams{
		Login:    user.Login,
		Password: user.Password,
		Username: user.Username,
	})
	if isConstraintViolation(err, uniqueLoginConstraint) {
		return 0, domain.ErrLoginAlreadyExists
	}
	if isConstraintViolation(err, uniqueUsernameConstraint) {
		return 0, domain.ErrUsernameAlreadyExists
	}
	if err != nil {
		return 0, fmt.Errorf("a.db.InsertUser: %w", err)
	}
	return domain.UserID(id), nil
}

func (a *Adapter) GetUserByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	user, err := a.db.GetUserByID(ctx, int64(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("a.db.GetUserByID: %w", err)
	}
	return convertUserFromDB(user), nil
}

func (a *Adapter) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	user, err := a.db.GetUserByLogin(ctx, login)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("a.db.GetUserByLogin: %w", err)
	}
	return convertUserFromDB(user), nil
}

func convertUserFromDB(user postgres_client.User) domain.User {
	return domain.User{
		ID:       domain.UserID(user.ID),
		ImageID:  domain.ImageID(user.ImageID),
		Login:    user.Login,
		Password: user.Password,
	}
}

func isConstraintViolation(err error, constraintName string) bool {
	if err == nil {
		return false
	}
	pgxErr := &pgconn.PgError{}

	if ok := errors.As(err, &pgxErr); ok {
		return pgxErr.ConstraintName == constraintName
	}

	return false
}
