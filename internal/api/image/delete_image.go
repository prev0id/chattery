package image

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	avatar, err := s.image.DeleteUserImage(ctx, userID)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertDeleteImageResponse(avatar))
}
