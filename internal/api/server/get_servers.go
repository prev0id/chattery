package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) GetServers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	servers, err := s.server.GetUserServers(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertGetServersResponse(servers))
}
