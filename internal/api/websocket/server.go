package websocket

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"

	"chattery/internal/api"
	wsmanager "chattery/internal/service/websocket_manager"
)

var (
	wsOptions = &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	}
)

const (
	readMessageLimit = api.MiB
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	user userService
	ws   *wsmanager.WebsocketManager
}

func New(user userService, ws *wsmanager.WebsocketManager) *Server {
	return &Server{
		user: user,
		ws:   ws,
	}
}

func (*Server) Pattern() string {
	return "/ws"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.user.AuthRequiredMiddleware)

		withAuthRouter.Get("/", s.Websocket)
	})
}
