package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) ListTopicMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, requestErr := bind.JSON[ListTopicMessagesRequest](r)
	if requestErr != nil {
		render.Error(w, r, requestErr)
		return
	}

	cursor := convertTopicCursorRequest(request.Cursor)

	var (
		messages   []*domain.TopicMessage
		nextCursor *domain.TopicCursor
		err        error
	)

	if cursor != nil && cursor.MessageID != 0 {
		messages, nextCursor, err = s.server.NextPagesOfTopicMessages(ctx, cursor)
	} else {
		messages, nextCursor, err = s.server.FirstPageOfTopicMessages(ctx, cursor)
	}

	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertListTopicMessagesResponse(nextCursor, messages))
}
