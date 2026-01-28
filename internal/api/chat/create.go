package chat_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

// CreatePublic создает публичный чат
func (s *Server) CreatePublic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[CreatePublicChatRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	id, err := s.chat.CreatePublicChat(ctx, userID, request.Name)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, CreateChatResponse{ID: id.I64()})
}

// CreatePrivate создает личный чат между двумя пользователями
func (s *Server) CreatePrivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[CreatePrivateChatRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	id, err := s.chat.CreatePrivateChat(ctx, userID, domain.UserID(request.WithUserID))
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, CreateChatResponse{ID: id.I64()})
}
