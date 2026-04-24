package dm

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostMessageRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	message := convertPostMessageRequest(request, userID)

	if err := s.dm.CreateDMMessage(ctx, message); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
