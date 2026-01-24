package user_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

// Login аутентификация по логину и паролю, ставит сессионную куку
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[LoginRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	user, err := s.user.GetByCredentials(ctx, domain.Login(request.Login), request.Password)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.user.CreateSession(ctx, w, user.ID); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
