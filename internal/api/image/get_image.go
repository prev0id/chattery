package image

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	"chattery/internal/utils/identicon"
)

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	username := domain.Username(chi.URLParam(r, "username"))

	user, err := s.user.GetByUsername(username)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if user.AvatarID != "" {
		http.Error(w, "custom avatars not implemented yet", http.StatusNotImplemented)
		return
	}

	imgBytes, err := identicon.Bytes([]byte(username), 420)
	if err != nil {
		http.Error(w, "failed to generate image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=600")
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(imgBytes) // #nosec G705 false-positive
}
