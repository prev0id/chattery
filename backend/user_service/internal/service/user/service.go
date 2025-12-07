package user

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	postgres_adapter "chattery/backend/user_service/internal/adapter/postgres"
	"chattery/backend/user_service/internal/config"
	"chattery/backend/user_service/internal/domain"
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

func (s *Service) CreateSessionInvalidation() domain.Session {
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

func (s *Service) Validate(session string) (domain.Claims, error) {
	token, err := jwt.ParseWithClaims(session, &domain.Claims{}, func(t *jwt.Token) (any, error) {
		return s.key, nil
	})
	if err != nil {
		slog.Info("jwt parsing error", slog.String("error", err.Error()))
		return domain.Claims{}, domain.ErrInvalidSession
	}
	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		slog.Info("session doesn't match")
		return domain.Claims{}, domain.ErrInvalidSession
	}
	return *claims, nil
}

func (s *Service) Update(ctx context.Context, session string, updated domain.User) error {
	claims, err := s.Validate(session)
	if err != nil {
		return err
	}

	user, err := s.db.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("s.db.GetUserByID: %w", err)
	}

	if updated.ImageID != "" {
		user.ImageID = updated.ImageID
	}
	if updated.Login != "" {
		user.Login = updated.Login
	}
	if updated.Password != nil {
		user.Password = updated.Password
	}

	if err := s.db.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("s.db.UpdateUser: %w", err)
	}
	return nil
}

func (s *Service) createSessionForUser(user domain.User) (domain.Session, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, domain.Claims{
		UserID:   user.ID,
		Username: user.Username,
		ImageID:  user.ImageID,
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
