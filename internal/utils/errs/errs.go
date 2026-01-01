// Package errs provides typed application errors with HTTP semantics
package errs

import (
	"bytes"
	"chattery/internal/domain"
	"fmt"
	"net/http"
)

// Error is a domain-aware application error
type Error struct {
	// related user, if any
	User domain.Username `json:"user,omitempty"`
	// semantic error class
	Kind Kind `json:"kind,omitempty"`
	// internal debug messages
	Debug []Debug `json:"debug,omitempty"`
	// public message
	Message string `json:"message,omitempty"`
	// root cause
	Err error `json:"err,omitempty"`
}

// Debug is a developer diagnostic message
type Debug string

// Kind classifies application errors
type Kind int

// Message public message
type Message string

const (
	Unknown Kind = iota
	InvalidRequest
	Unauthorized
	Permission
	Exist
	NotFound
	Internal
)

// StatusCode maps Kind to HTTP status
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

// String returns a short human-readable description
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

// E builds an Error from mixed arguments
func E(args ...any) *Error {
	err := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case domain.Username:
			err.User = arg
		case Kind:
			err.Kind = arg
		case Debug:
			err.Debug = append(err.Debug, arg)
		case error:
			err.Err = arg
		case string:
			err.Message = arg
		}
	}
	return err
}

// Error implements the error interface.
func (err *Error) Error() string {
	if err.Message == "" {
		return err.Kind.String()
	}
	return err.Kind.String() + ": " + err.Message
}

// String returns a multi-line diagnostic message.
func (err *Error) String() string {
	builder := new(bytes.Buffer)

	if err.Message != "" {
		fmt.Fprintf(builder, "MESSAGE=%q\n", err.Kind.String())
	}

	if err.Kind != 0 {
		fmt.Fprintf(builder, "KIND=%q\n", err.Kind.String())
	}

	if err.User != "" {
		fmt.Fprintf(builder, "USER=%q\n", err.User.String())
	}

	if err.Err != nil {
		fmt.Fprintf(builder, "ERROR=%q\n", err.Err.Error())
	}

	if len(err.Debug) > 0 {
		builder.WriteString("DEBUG INFO: ")
		for idx, debugMsg := range err.Debug {
			fmt.Fprintf(builder, "%d. %q\n", idx+1, debugMsg)
		}
	}

	if builder.Len() == 0 {
		return "no error"
	}

	return builder.String()
}

func (err *Error) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// pad writes only when buffer is not empty
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

// Is checks error kind
func Is(err error, kind Kind) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Kind == kind
}

// D adds debug message to incoming error
func D(message string, incoming error) *Error {
	err := FromError(incoming)
	err.Debug = append(err.Debug, Debug(message))
	return err
}

// FromError converts a standard error to *Error
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if converted, ok := err.(*Error); ok {
		return converted
	}
	return E(err)
}
