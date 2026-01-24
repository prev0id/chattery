package user_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

const searchQueryName = "query"

// Search поиск профиля
func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	query := r.URL.Query().Get(searchQueryName)

	users, err := s.user.Search(ctx, userID, query)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, converSearchResponse(users))
}
