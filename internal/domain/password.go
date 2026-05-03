package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	plain  string
	hashed []byte
}

func NewPassword(plain string, login Email) Password {
	rawWithSalt := addSalt(plain, login)
	hashed, _ := bcrypt.GenerateFromPassword(rawWithSalt, bcrypt.DefaultCost)

	return Password{
		plain:  plain,
		hashed: hashed,
	}
}

func NewHashedPassword(hashed []byte) Password {
	return Password{
		hashed: append([]byte(nil), hashed...),
	}
}

func (pass Password) Plain() string {
	return pass.plain
}

func (pass Password) Hashed() []byte {
	return append([]byte(nil), pass.hashed...)
}

func (pass Password) Equal(plain string, login Email) bool {
	rawWithSalt := addSalt(plain, login)
	err := bcrypt.CompareHashAndPassword(pass.hashed, rawWithSalt)
	return err == nil
}

func addSalt(plain string, login Email) []byte {
	return fmt.Appendf(nil, "%s$%s", plain, login.String())
}
