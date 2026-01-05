package userapi

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

type UpdateRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := domain.UsernameFromContext(ctx)

	request, err := bind.Json[UpdateRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validateUpdateRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	updated := convertUpdateRequest(request)

	if err := s.user.UpdateUser(ctx, username, updated); err != nil {
		render.Error(w, r, err)
		return
	}

	if updated.Username != domain.UserUnknown {
		username = updated.Username
	}

	if err := s.user.CreateSession(ctx, w, updated.Username); err != nil {
		render.Error(w, r, err)
		return
	}
}

func validateUpdateRequest(req *UpdateRequest) error {
	if req.Username != "" {
		if err := validate.Username(req.Username); err != nil {
			return err
		}
	}

	if req.Password != "" {
		if err := validate.NotEmpty(req.Login, validate.LoginFieldName); err != nil {
			return err
		}

		if err := validate.Password(req.Password); err != nil {
			return err
		}
	}

	if req.Login != "" {
		if err := validate.Login(req.Login); err != nil {
			return err
		}
	}
	return nil
}

func convertUpdateRequest(req *UpdateRequest) *domain.User {
	login := domain.Login(req.Login)
	return &domain.User{
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}
