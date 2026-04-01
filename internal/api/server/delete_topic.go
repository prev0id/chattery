package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[DeleteTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.DeleteTopic(ctx, domain.TopicID(request.TopicID), userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
