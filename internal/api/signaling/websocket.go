package signalingapi

import (
	"net/http"

	"github.com/coder/websocket"

	"chattery/internal/domain"
	"chattery/internal/service/subscription"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

func (s *Server) WebsocketEntrypoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)
	session := domain.GetSessionFromRequest(r)

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

	ctx, subscriber := subscription.New(ctx, userID, session, conn)

	s.chat.Register(ctx, subscriber)
	// s.webrtc.Register(ctx, user, subscriber)

	subscriber.Reader(ctx)
}
