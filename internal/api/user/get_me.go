package user_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	user, err := s.user.GetByID(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertGetMeResponse(user))
}
