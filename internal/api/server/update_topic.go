package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	request, err := bind.JSON[UpdateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	topic := convertUpdateTopicRequest(request)

	err = s.server.UpdateTopic(ctx, topic, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
