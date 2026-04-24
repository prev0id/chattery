package server

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostJoinServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostJoinServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	participant := convertPostJoinServerRequest(request, userID)

	if err = s.server.AddParticipant(ctx, participant); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
