package domain

import (
	"context"
)

const UserIsUnknown UserID = 0

type UserID int64

func (id UserID) I64() int64 { return int64(id) }

type userIDContextKeyType struct{}

var userIDContextKey userIDContextKeyType

func UserIDFromContext(ctx context.Context) UserID {
	id, ok := ctx.Value(userIDContextKey).(UserID)
	if !ok {
		return UserIsUnknown
	}
	return id
}

func UserIDToContext(ctx context.Context, user UserID) context.Context {
	return context.WithValue(ctx, userIDContextKey, user)
}
