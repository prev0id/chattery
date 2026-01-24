package chat_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

// Join войти в чат
func (s *Server) Join(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[JoinRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.chat.JoinChat(ctx, userID, domain.ChatID(request.ID)); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
