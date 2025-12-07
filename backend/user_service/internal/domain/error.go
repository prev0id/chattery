package domain

import (
	"errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrLoginAlreadyExists    = errors.New("login already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPasswordsDontMatch    = errors.New("passwords don't match")
	ErrInvalidSession        = errors.New("invalid session")
)

var ErrToGRPCCode = map[error]codes.Code{
	ErrUserNotFound:          codes.NotFound,
	ErrLoginAlreadyExists:    codes.InvalidArgument,
	ErrUsernameAlreadyExists: codes.InvalidArgument,
	ErrPasswordsDontMatch:    codes.PermissionDenied,
	ErrInvalidSession:        codes.Unauthenticated,
}
