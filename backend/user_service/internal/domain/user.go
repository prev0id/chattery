package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltPattern = "%s$%s"
	cost        = bcrypt.DefaultCost
)

type (
	UserID  int64
	ImageID string
)

type User struct {
	ID          UserID
	ImageID     ImageID
	Login       string
	Password    []byte
	Username    string
	rawPassword string
}

func (u *User) SetRawPassword(password string) {
	u.rawPassword = password
	withSalt := addSalt(password, u.Login)
	encrypted, _ := bcrypt.GenerateFromPassword(withSalt, cost)
	u.Password = encrypted
}

func (u *User) PasswordEqual(password string) bool {
	withSalt := addSalt(password, u.Login)
	err := bcrypt.CompareHashAndPassword(u.Password, withSalt)
	return err == nil
}

func addSalt(password string, email string) []byte {
	return fmt.Appendf(nil, saltPattern, email, password)
}
