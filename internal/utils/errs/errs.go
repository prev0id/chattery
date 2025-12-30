package errs

import (
	"bytes"
	"chattery/internal/domain"
	"errors"
	"net/http"
)

type Error struct {
	User domain.Username
	Kind Kind
	Err  error
}

type Kind int

const (
	Unknown Kind = iota
	InvalidRequest
	Unauthorized
	Permission
	Exist
	NotFound
	Internal
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
	case InvalidRequest:
		return "invalid request"
	case Unauthorized:
		return "unauthorized"
	case Permission:
		return "forbidden"
	case Exist:
		return "item already exists"
	case NotFound:
		return "item not found"
	default:
		return "unknown error"
	}
}

func E(args ...any) *Error {
	err := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case domain.Username:
			err.User = arg
		case Kind:
			err.Kind = arg
		case error:
			err.Err = arg
		case string:
			err.Err = errors.New(arg)
		}
	}
	return err
}

func (err *Error) Error() string {
	builder := new(bytes.Buffer)

	if err.User != "" {
		builder.WriteString("user ")
		builder.WriteString(string(err.User))
	}
	if err.Kind != 0 {
		pad(builder, ": ")
		builder.WriteString(err.Kind.String())
	}
	if err.Err != nil {
		pad(builder, ": ")
		builder.WriteString(err.Err.Error())
	}
	if builder.Len() == 0 {
		return "no error"
	}
	return builder.String()
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.Kind == kind
}
