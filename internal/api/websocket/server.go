package websocket_api

import (
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/service/websocket_manager"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type Server struct {
	userService userService
	wsManager   *websocket_manager.WebsocketManager
}

func New(user userService, wsManager *websocket_manager.WebsocketManager) *Server {
	return &Server{
		userService: user,
		wsManager:   wsManager,
	}
}

func (s *Server) Pattern() string {
	return "/v1/ws"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.userService.AuthRequiredMiddleware)

		withAuthRouter.Get("/dm/{dm_id}", s.DMWebsocket)
		withAuthRouter.Get("/text_topic/{topic_id}", s.TextTopicWebsocket)
	})
}

func (s *Server) establishWebsocket(w http.ResponseWriter, r *http.Request, userID domain.UserID, channelType websocket_manager.ChannelType, channelID int64) {
	ctx := r.Context()

	options := &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	}

	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		err = errors.E(err).
			Kind(errors.InvalidRequest).
			Message("unable to upgrade to websocket").
			Debug("websocket.Accept")
		render.Error(w, r, err)
		return
	}

	connection := websocket_manager.NewConnection(s.wsManager, userID, channelType, channelID, conn, ctx)
	s.wsManager.RegisterConnection(connection)

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
