package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
	"chattery/internal/utils/validate"
)

func (s *Server) PostTopicCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostCreateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if err = validatePostCreateTopic(request); err != nil {
		render.Error(w, r, err)
		return
	}

	topic := convertPostCreateTopicRequest(request)

	topicID, err := s.server.CreateTopic(ctx, topic, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertPostCreateTopicResponse(topicID))
}

func validatePostCreateTopic(request *PostCreateTopicRequest) error {
	if err := validate.TopicName(request.Name); err != nil {
		return err
	}

	return validate.TopicType(request.Type)
}
