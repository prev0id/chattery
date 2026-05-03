package domain

import (
	"context"
	"crypto/rand"
	"net/http"
)

const (
	NoSession Session = ""

	SessionCookieName = "__Session"
)

type Session string

type sessionErrorContextKeyType struct{}

var sessionErrorContextKey sessionErrorContextKeyType

func NewSession() Session {
	return Session(rand.Text())
}

func (s Session) Valid() bool {
	value := s.String()
	if len(value) < 26 || len(value) > 256 {
		return false
	}

	for _, r := range value {
		if r <= ' ' || r >= 127 {
			return false
		}
	}

	return true
}

func (s Session) String() string {
	return string(s)
}

func GetSessionFromRequest(r *http.Request) Session {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return NoSession
	}

	return Session(cookie.Value)
}

func SessionErrorFromContext(ctx context.Context) error {
	err, ok := ctx.Value(sessionErrorContextKey).(error)
	if !ok {
		return nil
	}
	return err
}

func SessionErrorToContext(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, sessionErrorContextKey, err)
}
