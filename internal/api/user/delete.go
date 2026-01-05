package user

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := domain.UsernameFromContext(ctx)
	session := domain.GetSessionFromRequest(r)

	if err := s.user.DeleteUser(ctx, username); err != nil {
		render.Error(w, r, errors.E(err).Debug("s.user.DeleteUser"))
		return
	}

	s.user.ClearSession(ctx, w, session)
	w.WriteHeader(http.StatusOK)
}
