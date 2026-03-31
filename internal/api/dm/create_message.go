package dm_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[CreateMessageRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	message := convertCreateMessageRequest(request, userID)

	if err := s.dm.CreateDMMessage(ctx, message); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
