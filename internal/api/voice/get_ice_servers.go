package voice

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) GetICEServers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	iceServers, err := s.voice.ICEServers(userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertGetICEServersResponse(iceServers))
}
