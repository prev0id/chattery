package chatapi

import (
	"chattery/internal/service/session"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
	"net/http"

	"github.com/coder/websocket"
)

func (s *Server) WebsocketEntrypoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := session.UsernameFromContext(ctx)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		err = errors.E(err).
			Kind(errors.InvalidRequest).
			Message("unable to upgrade to websocket").
			Debug("websocket.Accept")
		render.Error(w, r, err)
		return
	}

	defer conn.CloseNow()

	ctx, subscriber := s.signaling.Subscribe(ctx, username, conn)
	defer s.signaling.Unsubscribe(subscriber)

	s.chat.Register(ctx, subscriber)
	// s.webrtc.Register(ctx, user, subscriber)

	go subscriber.Writer(ctx)
	subscriber.Reader(ctx)
}
