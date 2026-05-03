package server

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostTopicMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, requestErr := bind.JSON[PostTopicMessagesRequest](r)
	if requestErr != nil {
		render.Error(w, r, requestErr)
		return
	}

	var (
		messages   []*domain.TopicMessage
		nextCursor *domain.TopicCursor
		err        error

		cursor = convertTopicCursorRequest(request.Cursor)
	)

	if cursor != nil && cursor.MessageID != 0 {
		messages, nextCursor, err = s.textTopic.NextPageOfMessages(ctx, cursor, userID)
	} else {
		messages, nextCursor, err = s.textTopic.FirstPageOfMessages(ctx, cursor, userID)
	}

	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertPostTopicMessagesResponse(nextCursor, messages, s.cache.ListByID()))
}
