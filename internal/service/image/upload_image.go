package image

import (
	"bytes"
	"context"
	stdimage "image"
	_ "image/gif" // decode
	"image/jpeg"
	_ "image/png" // decode
	"io"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

func (s *Service) UploadUserImage(ctx context.Context, userID domain.UserID, source io.Reader, size int64) (string, error) {
	if size <= 0 {
		return "", errutil.E().Kind(errutil.InvalidRequest).Message("image is empty")
	}
	if size > s.maxFileSize {
		return "", errutil.E().Kind(errutil.InvalidRequest).Message("image file is too large")
	}

	raw, err := readLimited(source, s.maxFileSize)
	if err != nil {
		return "", err
	}

	imgBytes, err := s.convertToJPEG(raw)
	if err != nil {
		return "", err
	}

	user, err := s.db.UserByID(ctx, userID)
	if err != nil {
		return "", errutil.E(err).Debug("s.db.UserByID")
	}

	imageID := domain.ImageID(objectName(user.Username))
	if err := s.storage.PutJPEGImage(ctx, imageID, imgBytes); err != nil {
		return "", errutil.E(err).Debug("s.storage.PutJPEGImage")
	}

	user.AvatarID = imageID
	if err := s.db.UpdateUser(ctx, user); err != nil {
		return "", errutil.E(err).Debug("s.db.UpdateUser")
	}

	return render.AvatarURL(user.Username), nil
}

func readLimited(source io.Reader, maxSize int64) ([]byte, error) {
	raw, err := io.ReadAll(io.LimitReader(source, maxSize+1))
	if err != nil {
		return nil, errutil.E(err).Debug("io.ReadAll")
	}
	if int64(len(raw)) > maxSize {
		return nil, errutil.E().Kind(errutil.InvalidRequest).Message("image file is too large")
	}

	return raw, nil
}

func (s *Service) convertToJPEG(raw []byte) ([]byte, error) {
	cfg, _, err := stdimage.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return nil, errutil.E(err).Kind(errutil.InvalidRequest).Message("invalid image")
	}
	if cfg.Width > s.maxWidth || cfg.Height > s.maxHeight {
		return nil, errutil.E().
			Kind(errutil.InvalidRequest).
			Messagef("image dimensions must not exceed %dx%d", s.maxWidth, s.maxHeight)
	}

	img, _, err := stdimage.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, errutil.E(err).Kind(errutil.InvalidRequest).Message("invalid image")
	}

	var b bytes.Buffer
	if err := jpeg.Encode(&b, img, &jpeg.Options{Quality: s.jpegQuality}); err != nil {
		return nil, errutil.E(err).Debug("jpeg.Encode")
	}

	return b.Bytes(), nil
}
