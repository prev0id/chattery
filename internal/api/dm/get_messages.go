package dm_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[GetMessagesRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	cursor := convertGetMessagesRequest(request)

	var (
		messages   []*domain.DMMessage
		nextCursor *domain.DMCursor
	)

	if cursor.MessageID == 0 {
		messages, nextCursor, err = s.dm.FirstPageOfDMMessages(ctx, userID, cursor)
	} else {
		messages, nextCursor, err = s.dm.NextPagesOfDMMessages(ctx, userID, cursor)
	}

	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertGetMessagesResponse(nextCursor, messages, s.cache.List()))
}
