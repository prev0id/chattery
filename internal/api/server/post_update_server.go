package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

func (s *Server) PostUpdateServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostUpdateServerRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validatePostUpdateServer(request); err != nil {
		render.Error(w, r, err)
	}

	err = s.server.UpdateServer(ctx, domain.ServerID(request.ServerID), request.Name, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validatePostUpdateServer(request *PostUpdateServerRequest) error {
	if err := validate.ServerName(request.Name); err != nil {
		return err
	}

	return nil
}
