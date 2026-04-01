package user_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

func (s *Server) PostUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostUpdateUserRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validatePostUpdateUserRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	updated := convertPostUpdateUserRequest(request, userID)

	if err := s.user.UpdateUser(ctx, updated); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validatePostUpdateUserRequest(req *PostUpdateUserRequest) error {
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
