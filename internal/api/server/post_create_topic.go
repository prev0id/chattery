package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostTopicCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostCreateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	topic := convertPostCreateTopicRequest(request)

	topicID, err := s.server.CreateTopic(ctx, topic, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertPostCreateTopicResponse(topicID))
}
