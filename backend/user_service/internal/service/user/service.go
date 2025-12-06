package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	postgres_adapter "chattery/backend/user_service/internal/adapter/postgres"
	"chattery/backend/user_service/internal/config"
	"chattery/backend/user_service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

const (
	sessionName    = "__Session"
	claimsUserID   = "user_id"
	claimsUsername = "username"
)

type Service struct {
	db         *postgres_adapter.Adapter
	expiration time.Duration
	key        []byte
}

func New(db *postgres_adapter.Adapter) *Service {
	return &Service{
		db:         db,
		expiration: config.Expiration,
		key:        []byte(config.SigningKey),
	}
}

func (s *Service) CreateUser(ctx context.Context, user domain.User) (domain.Session, error) {
	id, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return domain.Session{}, fmt.Errorf("s.db.CreateUser: %w", err)
	}
	user.ID = id

	session, err := s.createSessionForUser(user)
	if err != nil {
		return domain.Session{}, fmt.Errorf("s.createSessionForUser: %w", err)
	}
	return session, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (domain.Session, error) {
	user, err := s.db.GetUserByLogin(ctx, login)
	if err != nil {
		return domain.Session{}, fmt.Errorf("s.db.GetUserByLogin: %w", err)
	}

	if !user.PasswordEqual(password) {
		return domain.Session{}, domain.ErrPasswordsDontMatch
	}

	session, err := s.createSessionForUser(user)
	if err != nil {
		return domain.Session{}, fmt.Errorf("s.createSessionForUser: %w", err)
	}
	return session, nil
}

func (s *Service) createSessionForUser(user domain.User) (domain.Session, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, domain.Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	ss, err := token.SignedString(s.key)
	if err != nil {
		return domain.Session{}, fmt.Errorf("token.SignedString: %w", err)
	}

	return domain.Session{
		Name:     sessionName,
		Value:    ss,
		Path:     "/",
		MaxAge:   int(s.expiration.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, nil
}

func createSessionInvalidation() domain.Session {
	return domain.Session{
		Name:     sessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
