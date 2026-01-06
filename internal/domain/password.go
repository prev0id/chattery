package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type (
	Login    string
	Password []byte
)

func (l Login) String() string { return string(l) }

func NewPassword(raw string, login Login) Password {
	rawWithSalt := addSalt(raw, login)
	encrypted, _ := bcrypt.GenerateFromPassword(rawWithSalt, bcrypt.DefaultCost)
	return Password(encrypted)
}

func (pass Password) Equal(raw string, login Login) bool {
	rawWithSalt := addSalt(raw, login)
	err := bcrypt.CompareHashAndPassword(pass, rawWithSalt)
	return err == nil
}

func addSalt(raw string, login Login) []byte {
	return fmt.Appendf(nil, "%s$%s", raw, login)
}
