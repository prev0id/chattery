package domain

import (
	"crypto/rand"
)

type User struct {
	Username Username
	ImageID  ImageID
	Login    string
	Password Password
}

type Username string

func (u Username) String() string { return string(u) }

const UserUnknown Username = ""

type ImageID string

type Login string

type Session string

func NewSession() Session {
	return Session(rand.Text())
}

func (s Session) String() string {
	return string(s)
}
