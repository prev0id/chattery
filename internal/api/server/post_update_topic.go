package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) PostTopicUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[PostUpdateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	topic := convertPostUpdateTopicRequest(request)

	err = s.server.UpdateTopic(ctx, topic, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
