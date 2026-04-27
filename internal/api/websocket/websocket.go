package websocket

import (
	"net/http"
	"sync"

	"github.com/coder/websocket"

	"chattery/internal/domain"
)

func (s *Server) Websocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	conn, err := websocket.Accept(w, r, wsOptions)
	if err != nil {
		return
	}

	connection := s.ws.NewConnection(ctx, userID, conn)

	var wg sync.WaitGroup
	wg.Go(func() { connection.WritePump() })
	wg.Go(func() { connection.ReadPump() })
	wg.Wait()
}
