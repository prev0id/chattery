package image_api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chattery/internal/domain"
	identicon "chattery/internal/utils/identicon"
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

	w.Header().Set("Content-Type", "image/png")
	w.Write(imgBytes)
}
