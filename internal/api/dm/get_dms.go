package dm_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) GetDMs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := domain.UserIDFromContext(ctx)

	dms, err := s.dm.UserDMs(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertGetDMsResponse(dms, s.cache.List()))
}
