package user

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userService interface {
	GetByCredentials(ctx context.Context, login domain.Email, rawPassword string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error)
	UpdateUser(ctx context.Context, updated *domain.User) error
	DeleteUser(ctx context.Context, userID domain.UserID) error
	GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
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

func (*Server) Pattern() string {
	return "/v1/user"
}

func (s *Server) Route(router chi.Router) {
	router.Post("/create", s.PostCreateUser)
	router.Post("/login", s.PostLogin)

	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)
		withAuthRouter.Post("/logout", s.PostLogout)
		withAuthRouter.Put("/update", s.PostUpdateUser)
		withAuthRouter.Delete("/delete", s.DeleteMe)
		withAuthRouter.Get("/me", s.GetMe)
	})
}
