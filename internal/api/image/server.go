package image

import (
	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userStore interface {
	GetByUsername(username domain.Username) (*domain.User, error)
}

type Server struct {
	user userStore
}

func New(user userStore) *Server {
	return &Server{
		user: user,
	}
}

func (*Server) Pattern() string {
	return "/v1/image"
}

func (s *Server) Route(router chi.Router) {
	router.Get("/{username}.png", s.GetImage)
}
