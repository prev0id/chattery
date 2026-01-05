package user

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

type CreateRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.Json[CreateRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validateCreateRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	user := convertCreateRequest(request)

	if err := s.user.CreateUser(ctx, user); err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.user.CreateSession(ctx, w, user.Username); err != nil {
		render.Error(w, r, err)
		return
	}
}

func validateCreateRequest(req *CreateRequest) error {

	return nil
}

func convertCreateRequest(req *CreateRequest) *domain.User {
	login := domain.Login(req.Login)
	return &domain.User{
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}
