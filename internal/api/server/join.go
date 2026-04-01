package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) Join(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[JoinServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	participant := convertJoinServerRequest(request, userID)

	if err = s.server.AddParticipant(ctx, participant); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
