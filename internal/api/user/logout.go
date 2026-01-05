package userapi

import (
	"net/http"
)

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s.user.ClearSession(ctx, w, r)

	w.WriteHeader(http.StatusOK)
}
