package voice

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type voiceService interface {
	ICEServers(userID domain.UserID) ([]domain.VoiceICEServer, error)
}

type Server struct {
	user  userService
	voice voiceService
}

func New(user userService, voice voiceService) *Server {
	return &Server{
		user:  user,
		voice: voice,
	}
}

func (*Server) Pattern() string {
	return "/v1/voice"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Get("/ice-servers", s.GetICEServers)
	})
}
