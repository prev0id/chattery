package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[DeleteTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.DeleteTopic(ctx, domain.TopicID(request.TopicID))
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
