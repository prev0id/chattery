package user_api

import (
	"net/http"

	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

// Create создает новый профиль, ставит сессионную куку
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[CreateRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err := validateCreateRequest(request); err != nil {
		render.Error(w, r, err)
		return
	}

	user := convertCreateRequest(request)

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

func validateCreateRequest(req *CreateRequest) error {
	if err := validate.Username(req.Username); err != nil {
		return err
	}
	if err := validate.Password(req.Password); err != nil {
		return err
	}
	if err := validate.Login(req.Login); err != nil {
		return err
	}
	return nil
}
