package chat_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
	"chattery/internal/utils/sliceutil"
)

// ListPrivate получить список приватных чатов пользователя.
// В ответе для каждого чата добавляется последнее сообщение (если есть).
func (s *Server) ListPrivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	previews, err := s.chat.UserPrivateChats(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, PrivateChatsResponse{
		Chats: sliceutil.Map(previews, convertChatPreview),
	})
}

// ListPublic получить список приватных чатов пользователя.
// В ответе для каждого чата добавляется последнее сообщение (если есть).
func (s *Server) ListPublic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	previews, err := s.chat.UserPublicChats(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, PublicChatsResponse{
		Chats: sliceutil.Map(previews, convertChatPreview),
	})
}
