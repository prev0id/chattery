package user

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
	if err := validateUsername(req.Username); err != nil {
		return err
	}

	if err := validatePassword(req.Password, req.Login); err != nil {
		return err
	}

	return validateLogin(req.Login)
}

func validateUsername(username string) error {
	if username == "" {
		return nil
	}
	return validate.Username(username)
}

func validatePassword(password, login string) error {
	if password == "" {
		return nil
	}
	if err := validate.NotEmpty(login, validate.LoginFieldName); err != nil {
		return err
	}
	return validate.Password(password)
}

func validateLogin(login string) error {
	if login == "" {
		return nil
	}
	return validate.Login(login)
}
