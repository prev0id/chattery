package user_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

// DeleteMe полное удаление профиля
func (s *Server) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	if err := s.user.DeleteUser(ctx, userID); err != nil {
		render.Error(w, r, errors.E(err).Debug("s.user.DeleteUser"))
		return
	}

	s.user.ClearSession(ctx, w, r)
	w.WriteHeader(http.StatusOK)
}
