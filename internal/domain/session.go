package domain

import (
	"crypto/rand"
	"net/http"
)

const (
	NoSession Session = ""

	SessionCookieName = "__Session"
)

type Session string

func NewSession() Session {
	return Session(rand.Text())
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
