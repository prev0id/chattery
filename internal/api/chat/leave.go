package chat_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

// Leave выйти из чата
func (s *Server) Leave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[LeaveRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.chat.LeaveChat(ctx, userID, domain.ChatID(request.ID)); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
