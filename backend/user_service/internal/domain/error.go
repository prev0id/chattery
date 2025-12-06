package domain

import (
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrLoginAlreadyExists    = errors.New("login already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

var errToGRPCCode = map[error]codes.Code{
	ErrUserNotFound:          codes.NotFound,
	ErrLoginAlreadyExists:    codes.InvalidArgument,
	ErrUsernameAlreadyExists: codes.InvalidArgument,
}

func HandleGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if code, found := unwrapCode(err); found {
		return status.Error(code, err.Error())
	}
	slog.Error("[error_util] unknown error", slog.String("error", err.Error()))
	return err
}

func unwrapCode(err error) (codes.Code, bool) {
	for targetErr, code := range errToGRPCCode {
		if errors.Is(err, targetErr) {
			return code, true
		}
	}
	return codes.Internal, false
}
