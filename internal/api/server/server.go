package server_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type serverService interface {
}

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	user   userService
	server serverService
}

func New(user userService, server serverService) *Server {
	return &Server{
		user:   user,
		server: server,
	}
}

func (s *Server) Pattern() string {
	return "/v1/server"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		// GET /v1/server/list // with body
		// POST /v1/server/join
		// POST /v1/server/leave
		// POST /v1/server/update
		//
		// POST /v1/server/topic/create
		// POST /v1/server/topic/update
	})
}
