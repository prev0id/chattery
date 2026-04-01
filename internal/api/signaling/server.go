package signaling_api

import (
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/hub"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	userService userService
	hub         *hub.Hub
}

func New(user userService, hubInstance *hub.Hub) *Server {
	return &Server{
		userService: user,
		hub:         hubInstance,
	}
}

func (s *Server) Pattern() string {
	return "/v1/signaling"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.userService.AuthRequiredMiddleware)

		withAuthRouter.Get("/ws", s.WebsocketEntrypoint)
	})
}

func (s *Server) WebsocketEntrypoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		err = errors.E(err).
			Kind(errors.InvalidRequest).
			Message("unable to upgrade to websocket").
			Debug("websocket.Accept")
		render.Error(w, r, err)
		return
	}

	connection := hub.NewConnection(s.hub, userID, conn)
	s.hub.RegisterConnection(connection)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		connection.WritePump(ctx)
	}()

	go func() {
		defer wg.Done()
		connection.ReadPump(ctx)
	}()

	wg.Wait()
}
