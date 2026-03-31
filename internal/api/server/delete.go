package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[DeleteServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.DeleteServer(ctx, domain.ServerID(request.ServerID))
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
