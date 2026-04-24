package user

import (
	"net/http"

	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

// PostCreateUser создает новый профиль, ставит сессионную куку
func (s *Server) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[PostCreateUserRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err = validatePostCreateUserRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	user := convertPostCreateUserRequest(request)

	userID, err := s.user.CreateUser(ctx, user)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := s.user.CreateSession(ctx, w, userID); err != nil {
		render.Error(w, r, err)
		return
	}
}

func validatePostCreateUserRequest(req *PostCreateUserRequest) error {
	if err := validate.Username(req.Username); err != nil {
		return err
	}
	if err := validate.Password(req.Password); err != nil {
		return err
	}
	return validate.Login(req.Login)
}
