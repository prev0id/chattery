package websocket_api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/service/websocket_manager"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/render"
)

func (s *Server) TextTopicWebsocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	topicIDStr := chi.URLParam(r, "topic_id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		err = errors.E(err).
			Kind(errors.InvalidRequest).
			Message("invalid topic id").
			Debug("domain.ParseTopicID")
		render.Error(w, r, err)
		return
	}

	if err := s.wsManager.UserHasAccessToTextTopic(ctx, userID, domain.TopicID(topicID)); err != nil {
		err = errors.E(err).
			Kind(errors.InvalidRequest).
			Message("no access to text topic").
			Debug("s.wsManager.UserHasAccessToTextTopic")
		render.Error(w, r, err)
		return
	}

	s.establishWebsocket(w, r, userID, websocket_manager.ChannelTextTopic, int64(topicID))
}
