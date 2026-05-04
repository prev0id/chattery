package image

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/config"
	"chattery/internal/domain"
)

type userStore interface {
	GetByUsername(username domain.Username) (*domain.User, error)
}

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type imageService interface {
	GetUserImage(ctx context.Context, user *domain.User) ([]byte, error)
	UploadUserImage(ctx context.Context, userID domain.UserID, source io.Reader, size int64) (string, error)
	DeleteUserImage(ctx context.Context, userID domain.UserID) (string, error)
}

type Server struct {
	user              userService
	cache             userStore
	image             imageService
	maxUploadBodySize int64
}

func New(user userService, cache userStore, image imageService, cfg *config.Config) *Server {
	return &Server{
		user:              user,
		cache:             cache,
		image:             image,
		maxUploadBodySize: cfg.Image.MaxFileSize + 1024*1024,
	}
}

func (*Server) Pattern() string {
	return "/v1/image"
}

func (s *Server) Route(router chi.Router) {
	router.Get("/{username}.jpeg", s.GetImage)

	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)
		withAuthRouter.Post("/upload", s.PostUploadImage)
		withAuthRouter.Delete("/delete", s.DeleteImage)
	})
}
