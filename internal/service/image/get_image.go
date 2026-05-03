package image

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) GetUserImage(ctx context.Context, user *domain.User) ([]byte, error) {
	imageID := domain.ImageID(objectName(user.Username))

	exists, err := s.storage.ImageExists(ctx, imageID)
	if err != nil {
		return nil, errutil.E(err).Debug("s.storage.ImageExists")
	}
	if exists {
		return s.storage.GetImage(ctx, imageID)
	}

	imgBytes, err := s.identiconBytes(user.Username)
	if err != nil {
		return nil, errutil.E(err).Debug("s.identiconBytes")
	}

	if err := s.storage.PutJPEGImage(ctx, imageID, imgBytes); err != nil {
		return nil, errutil.E(err).Debug("s.storage.PutJPEGImage")
	}

	updated := *user
	updated.AvatarID = imageID
	if err := s.db.UpdateUser(ctx, &updated); err != nil {
		return nil, errutil.E(err).Debug("s.db.UpdateUser")
	}

	return imgBytes, nil
}
