package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type (
	Email    string
	Password []byte
)

func (l Email) String() string { return string(l) }

func NewPassword(raw string, login Email) Password {
	rawWithSalt := addSalt(raw, login)
	encrypted, _ := bcrypt.GenerateFromPassword(rawWithSalt, bcrypt.DefaultCost)
	return Password(encrypted)
}

func (pass Password) Equal(raw string, login Email) bool {
	rawWithSalt := addSalt(raw, login)
	err := bcrypt.CompareHashAndPassword(pass, rawWithSalt)
	return err == nil
}

func addSalt(raw string, login Email) []byte {
	return fmt.Appendf(nil, "%s$%s", raw, login)
}
