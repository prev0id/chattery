package domain

import (
	"context"
)

const UserUnknown Username = ""

type Username string

func (u Username) String() string {
	return string(u)
}

type usernameContextKeyType struct{}

var usernameContextKey usernameContextKeyType

func UsernameFromContext(ctx context.Context) Username {
	username, ok := ctx.Value(usernameContextKey).(Username)
	if !ok {
		return UserUnknown
	}
	return username
}

func UsernameToContext(ctx context.Context, user Username) context.Context {
	return context.WithValue(ctx, usernameContextKey, user)
}
