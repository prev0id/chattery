package image

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := domain.Username(chi.URLParam(r, "username"))

	user, err := s.cache.GetByUsername(username)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	imgBytes, err := s.image.GetUserImage(ctx, user)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=600")
	w.Header().Set("Content-Type", "image/jpeg")
	_, _ = w.Write(imgBytes) // #nosec G705 false-positive
}
