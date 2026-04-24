package user

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
)

func convertUserFromDB(user *postgres.User) *domain.User {
	return &domain.User{
		ID:       domain.UserID(user.ID),
		Username: domain.Username(user.Username),
		AvatarID: domain.ImageID(user.AvatarID),
		Login:    domain.Email(user.Login),
		Password: user.Password,
	}
}
