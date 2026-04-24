package websocket //nolint:dupl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

func (s *Server) TextTopicWebsocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	topicIDStr := chi.URLParam(r, "topic_id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		err = errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("invalid topic id").
			Debug("domain.ParseTopicID")
		render.Error(w, r, err)
		return
	}

	if err := s.wsManager.UserHasAccessToTextTopic(ctx, userID, domain.TopicID(topicID)); err != nil {
		err = errutil.E(err).
			Kind(errutil.InvalidRequest).
			Message("no access to text topic").
			Debug("s.wsManager.UserHasAccessToTextTopic")
		render.Error(w, r, err)
		return
	}

	s.establishWebsocket(w, r, userID, domain.ChannelTextTopic, topicID)
}
