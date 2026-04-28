package dm

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostReadRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	dmID := domain.DMID(request.DMID)
	messageID := domain.DMMessageID(request.MessageID)

	if err := s.dm.MarkDMMessageRead(ctx, userID, dmID, messageID); err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
