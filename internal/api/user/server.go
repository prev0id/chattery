package userapi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userService interface {
	ValidateCredentials(ctx context.Context, login domain.Login, rawPassword string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, username domain.Username, updated *domain.User) error
	DeleteUser(ctx context.Context, username domain.Username) error
	CreateSession(ctx context.Context, w http.ResponseWriter, user domain.Username) error
	ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request)
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	user userService
}

func New(user userService) *Server {
	return &Server{
		user: user,
	}
}

func (s *Server) Pattern() string {
	return "/v1/user"
}

func (s *Server) Route(router chi.Router) {
	router.Post("/create", s.Create)
	router.Post("/login", s.Login)

	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Post("/logout", s.Logout)
		withAuthRouter.Put("/update", s.Update)
		withAuthRouter.Delete("/delete", s.Delete)
	})
}
