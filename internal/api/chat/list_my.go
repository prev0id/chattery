package chat_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

// ListMy получить список своих чатов
func (s *Server) ListMy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	chats, err := s.chat.UserChats(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, converMyChatsResponse(chats))
}
