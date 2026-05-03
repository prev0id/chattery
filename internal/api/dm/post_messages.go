package dm

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostMessagesRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	cursor := convertPostMessagesRequest(request)

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

	render.JSON(w, r, convertPostMessagesResponse(nextCursor, messages, s.cache.ListByID()))
}
