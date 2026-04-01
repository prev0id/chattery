package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) DeleteServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[DeleteServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.DeleteServer(ctx, domain.ServerID(request.ServerID), userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
