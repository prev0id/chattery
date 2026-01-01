package errors

import (
	"chattery/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Kind int8

const (
	Internal Kind = iota
	InvalidRequest
	Unauthorized
	Permission
	Exist
	NotFound
)

func (k Kind) StatusCode() int {
	switch k {
	case InvalidRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Permission:
		return http.StatusForbidden
	case Exist:
		return http.StatusConflict
	case NotFound:
		return http.StatusNotFound
	case Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (k Kind) String() string {
	switch k {
	case InvalidRequest:
		return "invalid request"
	case Unauthorized:
		return "unauthorized"
	case Permission:
		return "forbidden"
	case Exist:
		return "already exists"
	case NotFound:
		return "not found"
	default:
		return "unknown error"
	}
}

type Error struct {
	user    domain.Username
	kind    Kind
	debug   []string
	message string
	err     error
}

func E(errs ...error) *Error {
	if len(errs) == 0 {
		return &Error{}
	}

	err := errs[0]
	if domainErr, ok := err.(*Error); ok {
		return domainErr
	}

	return &Error{
		err: errs[0],
	}
}

func (err *Error) Error() string {
	if err.message == "" {
		return err.kind.String()
	}
	return err.kind.String() + ": " + err.message
}

func (e *Error) Kind(kind Kind) *Error {
	e.kind = kind
	return e
}

func (e *Error) GetKind() Kind {
	return e.kind
}

func (e *Error) User(user domain.Username) *Error {
	e.user = user
	return e
}

func (e *Error) GetUser() domain.Username {
	return e.user
}

func (e *Error) Debug(messages ...string) *Error {
	e.debug = append(e.debug, messages...)
	return e
}

func (e *Error) GetDebug() []string {
	return e.debug
}

func (e *Error) Message(message string) *Error {
	e.message = message
	return e
}

func (e *Error) Messagef(format string, args ...any) *Error {
	e.message = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) GetMessage() string {
	return e.message
}

func (e *Error) GetError() error {
	return e.err
}

func (e *Error) Log(message string, attr ...slog.Attr) *Error {
	attr = append(attr,
		slog.Group("error",
			slog.String("user", e.user.String()),
			slog.String("kind", e.kind.String()),
			slog.Any("debug", e.debug),
			slog.String("message", e.message),
			slog.Any("err", e.err),
		),
	)

	slog.LogAttrs(context.Background(), slog.LevelError, message, attr...)

	return e
}

func (e *Error) LogFatal() {
	e.Log("got unrecoverable error")
	os.Exit(1)
}
