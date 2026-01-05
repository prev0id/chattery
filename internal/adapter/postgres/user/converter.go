package useradapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
)

func convertUserFromDB(user *postgres.User) *domain.User {
	return &domain.User{
		Username: domain.Username(user.Username),
		AvatarID: domain.ImageID(user.AvatarID),
		Login:    domain.Login(user.Login),
		Password: user.Password,
	}
}
