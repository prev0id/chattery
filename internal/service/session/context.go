package session

import (
	"chattery/internal/domain"
	"context"
)

type usernameContextKeyType struct{}

var usernameContextKey usernameContextKeyType

func UsernameFromContext(ctx context.Context) domain.Username {
	username, ok := ctx.Value(usernameContextKey).(domain.Username)
	if !ok {
		return domain.UnknownUsername
	}
	return username
}

func usernameToContext(ctx context.Context, user domain.Username) context.Context {
	return ctx
}
