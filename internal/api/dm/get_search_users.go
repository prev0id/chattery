package dm

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

const searchQueryName = "query"

func (s *Server) SearchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)
	query := r.URL.Query().Get(searchQueryName)

	users, err := s.dm.SearchUsersWithoutDM(ctx, userID, query)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertSearchUsersResponse(users))
}
