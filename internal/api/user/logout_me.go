package user_api

import (
	"net/http"
)

// LogoutMe разлогин, удаляет сессионную куку
func (s *Server) LogoutMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s.user.ClearSession(ctx, w, r)

	w.WriteHeader(http.StatusOK)

	http.Redirect(w, r, "/login", http.StatusOK)
}
