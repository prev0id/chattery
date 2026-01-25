package signaling_api

import (
	"net/http"

	"github.com/coder/websocket"

	"chattery/internal/api/signaling/subscriber"
	"chattery/internal/domain"
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

	sub := subscriber.New(conn).
		WithSession(session).
		WithUserID(userID).
		WithChatService(s.chat)

	s.chat.Register(sub)
	defer s.chat.Unregister(sub)

	sub.Read(ctx)
}
