package websocket

import (
	"context"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

type userService interface {
	AuthRequiredMiddleware(next http.Handler) http.Handler
}

type websocketManager interface {
	UserHasAccessToDM(ctx context.Context, userID domain.UserID, dmID domain.DMID) error
	UserHasAccessToTextTopic(ctx context.Context, userID domain.UserID, topicID domain.TopicID) error
	NewConnection(ctx context.Context, userID domain.UserID, channelType domain.ChannelType, channelID int64, ws *websocket.Conn) domain.Connection
}

type Server struct {
	userService userService
	wsManager   websocketManager
}

func New(user userService, wsManager websocketManager) *Server {
	return &Server{
		userService: user,
		wsManager:   wsManager,
	}
}

func (*Server) Pattern() string {
	return "/ws"
}

func (s *Server) Route(router chi.Router) {
	router.Group(func(withAuthRouter chi.Router) {
		withAuthRouter.Use(s.userService.AuthRequiredMiddleware)

		withAuthRouter.Get("/dm/{dm_id}", s.DMWebsocket)
		withAuthRouter.Get("/text_topic/{topic_id}", s.TextTopicWebsocket)
	})
}

func (s *Server) establishWebsocket(w http.ResponseWriter, r *http.Request, userID domain.UserID, channelType domain.ChannelType, channelID int64) {
	ctx := r.Context()

	options := &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	}

	conn, err := websocket.Accept(w, r, options)
	if err != nil {
		err = errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("unable to upgrade to websocket").
			Debug("websocket.Accept")
		render.Error(w, r, err)
		return
	}

	connection := s.wsManager.NewConnection(ctx, userID, channelType, channelID, conn)

	var wg sync.WaitGroup

	wg.Go(func() {
		connection.WritePump(ctx)
	})

	wg.Go(func() {
		connection.ReadPump(ctx)
	})

	wg.Wait()
}
