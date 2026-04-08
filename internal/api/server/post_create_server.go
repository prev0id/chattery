package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

func (s *Server) PostCreateServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostCreateServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validatePostCreateServer(request); err != nil {

	}

	serverID, err := s.server.CreateServer(ctx, request.Name, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertPostCreateServerResponse(serverID))
}

func validatePostCreateServer(request *PostCreateServerRequest) error {
	if err := validate.ServerName(request.Name); err != nil {
		return err
	}

	return nil
}
