package websocket //nolint:dupl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

func (s *Server) DMWebsocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	dmIDStr := chi.URLParam(r, "dm_id")
	dmID, err := strconv.ParseInt(dmIDStr, 10, 64)
	if err != nil {
		err = errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("invalid dm_id").
			Debug("domain.ParseDMID")
		render.Error(w, r, err)
		return
	}

	if err := s.wsManager.UserHasAccessToDM(ctx, userID, domain.DMID(dmID)); err != nil {
		err = errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("no access to dm").
			Debug("s.wsManager.UserHasAccessToDM")
		render.Error(w, r, err)
		return
	}

	s.establishWebsocket(w, r, userID, domain.ChannelDM, dmID)
}
