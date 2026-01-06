package userapi

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	if err := s.user.DeleteUser(ctx, userID); err != nil {
		render.Error(w, r, errors.E(err).Debug("s.user.DeleteUser"))
		return
	}

	s.user.ClearSession(ctx, w, r)
	w.WriteHeader(http.StatusOK)
}
