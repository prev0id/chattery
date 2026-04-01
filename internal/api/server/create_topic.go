package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) CreateTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[CreateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	topic := convertCreateTopicRequest(request)

	topicID, err := s.server.CreateTopic(ctx, topic, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertTopicResponse(topicID))
}
