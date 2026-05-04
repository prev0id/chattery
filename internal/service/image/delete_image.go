package image

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/render"
)

func (s *Service) DeleteUserImage(ctx context.Context, userID domain.UserID) (string, error) {
	user, err := s.db.UserByID(ctx, userID)
	if err != nil {
		return "", errutil.E(err).Debug("s.db.UserByID")
	}

	imgBytes, err := s.identiconBytes(user.Username)
	if err != nil {
		return "", errutil.E(err).Debug("s.identiconBytes")
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
