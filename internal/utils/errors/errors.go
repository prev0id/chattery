package errors

import (
	"errors"
	"fmt"
	"net/http"
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
	default:
		return http.StatusInternalServerError
	}
}

func (k Kind) String() string {
	switch k {
	case Internal:
		return "internal error"
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
	err     error
	message string
	debug   []string
	kind    Kind
}

func E(errs ...error) *Error {
	if len(errs) == 0 {
		return &Error{}
	}

	err := errs[0]
	domainErr := &Error{}
	if errors.As(err, &domainErr) {
		return domainErr
	}

	return &Error{
		err: errs[0],
	}
}

func (e *Error) Error() string {
	if e.message == "" {
		return e.kind.String()
	}
	return e.kind.String() + ": " + e.message
}

func (e *Error) Kind(kind Kind) *Error {
	e.kind = kind
	return e
}

func (e *Error) GetKind() Kind {
	return e.kind
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

func Is(kind Kind, err error) bool {
	return E(err).kind == kind
}
