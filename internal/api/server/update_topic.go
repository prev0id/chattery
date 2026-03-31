package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[UpdateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	err = s.server.UpdateTopic(ctx, domain.TopicID(request.TopicID), request.Name)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
