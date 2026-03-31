package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) CreateTopicMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[CreateTopicMessageRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	messageID, err := s.server.CreateTopicMessage(ctx, domain.TopicID(request.TopicID), userID, request.Text)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertCreateTopicMessageResponse(messageID))
}
