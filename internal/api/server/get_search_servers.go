package server

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

const searchQueryName = "query"

func (s *Server) SearchServers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)
	query := r.URL.Query().Get(searchQueryName)

	servers, err := s.server.SearchServersNotJoined(ctx, userID, query)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertSearchServersResponse(servers))
}
