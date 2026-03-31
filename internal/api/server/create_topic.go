package server_api

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/render"
)

func (s *Server) CreateTopic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := bind.JSON[CreateTopicRequest](r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	topic, err := s.server.CreateTopic(ctx, domain.ServerID(request.ServerID), request.Name, request.Type)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.Json(w, r, convertTopicResponse(topic))
}
