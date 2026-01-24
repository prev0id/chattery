package user_api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userService interface {
	GetByCredentials(ctx context.Context, login domain.Login, rawPassword string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error)
	UpdateUser(ctx context.Context, updated *domain.User) error
	DeleteUser(ctx context.Context, userID domain.UserID) error
	Search(ctx context.Context, user domain.UserID, query string) ([]*domain.User, error)
	// TODO move to utils wrappers around http
	CreateSession(ctx context.Context, w http.ResponseWriter, userID domain.UserID) error
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

		withAuthRouter.Post("/me/logout", s.LogoutMe)
		withAuthRouter.Put("/me/update", s.UpdateMe)
		withAuthRouter.Delete("/me/delete", s.DeleteMe)

		withAuthRouter.Get("/info", s.Info)
		// withAuthRouter.Get("/search", s.Search)
		// withAuthRouter.Get("/list", s.List)
	})
}
