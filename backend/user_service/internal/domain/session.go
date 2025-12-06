package domain

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	Name     string
	Value    string
	Path     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

type Claims struct {
	jwt.RegisteredClaims

	UserID   UserID `json:"user_id"`
	Username string `json:"username"`
}
