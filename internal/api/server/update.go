package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[UpdateServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.UpdateServer(ctx, domain.ServerID(request.ServerID), request.Name)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
