package chatapi

import (
	"log/slog"
	"net/http"
)

func (s *Server) WebsocketEntrypoint(w http.ResponseWriter, r *http.Request) {
	slog.Info("hello world")
}
