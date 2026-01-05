package userapi

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.Json[LoginRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	user, err := s.user.ValidateCredentials(ctx, domain.Login(request.Login), request.Password)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.user.CreateSession(ctx, w, user.Username); err != nil {
		render.Error(w, r, err)
		return
	}
}
