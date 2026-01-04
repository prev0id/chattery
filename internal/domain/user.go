package domain

import (
	"crypto/rand"
)

type User struct {
	Username Username
	ImageID  ImageID
	Login    Login
	Password Password
}

type ImageID string

type Session string

func NewSession() Session {
	return Session(rand.Text())
}

func (s Session) String() string {
	return string(s)
}
