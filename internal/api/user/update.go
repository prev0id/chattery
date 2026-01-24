package user_api

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

// UpdateMe обновление учетных данных
func (s *Server) UpdateMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[UpdateRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validateUpdateRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	updated := convertUpdateRequest(request, userID)

	if err := s.user.UpdateUser(ctx, updated); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
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

func convertUpdateRequest(req *UpdateRequest, userID domain.UserID) *domain.User {
	login := domain.Login(req.Login)
	return &domain.User{
		ID:       userID,
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}
