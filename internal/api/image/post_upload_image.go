package image

import (
	"net/http"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

func (s *Server) PostUploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := domain.UserIDFromContext(ctx)

	r.Body = http.MaxBytesReader(w, r.Body, s.maxUploadBodySize)
	if err := r.ParseMultipartForm(s.maxUploadBodySize); err != nil {
		render.Error(w, r, errutil.E(err).Kind(errutil.InvalidRequest).Message("invalid multipart form"))
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		render.Error(w, r, errutil.E(err).Kind(errutil.InvalidRequest).Message("image file is required"))
		return
	}
	defer file.Close()

	avatar, err := s.image.UploadUserImage(ctx, userID, file, header.Size)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, r, convertPostUploadImageResponse(avatar))
}
