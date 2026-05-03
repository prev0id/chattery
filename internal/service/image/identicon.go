package image

import (
	"bytes"
	"image/jpeg"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) identiconBytes(username domain.Username) ([]byte, error) {
	identicon, err := s.identicon.Draw(username.String())
	if err != nil {
		return nil, errutil.E(err).Debug("s.identicon.Draw", username.String())
	}

	img := identicon.Image(identiconSize)

	var b bytes.Buffer
	if err := jpeg.Encode(&b, img, &jpeg.Options{Quality: s.jpegQuality}); err != nil {
		return nil, errutil.E(err).Debug("jpeg.Encode")
	}

	return b.Bytes(), nil
}
