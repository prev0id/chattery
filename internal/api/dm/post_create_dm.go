package dm_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostCreateDM(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostCreateDMRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	dmID, err := s.dm.CreateDM(ctx, userID, domain.UserID(request.ParticipantID))
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertPostCreateDMResponse(dmID))
}
