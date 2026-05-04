package user

import (
	"chattery/internal/domain"
	"chattery/internal/utils/render"
)

type PostCreateUserRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PostLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUsersResponse struct {
	Users []User
}

type PostUpdateUserRequest struct {
	Username        string `json:"username"`
	Login           string `json:"login"`
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
}

type GetMeResponse struct {
	Me User `json:"me"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	ID       int64  `json:"id"`
}

func convertPostCreateUserRequest(req *PostCreateUserRequest) *domain.User {
	login := domain.Email(req.Login)
	return &domain.User{
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}

func convertPostUpdateUserRequest(req *PostUpdateUserRequest, userID domain.UserID) *domain.User {
	login := domain.Email(req.Login)
	return &domain.User{
		ID:       userID,
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}

func convertUserResponse(user *domain.User) User {
	return User{
		ID:       user.ID.I64(),
		Username: user.Username.String(),
		Email:    user.Login.String(),
		Avatar:   render.AvatarURL(user.Username),
	}
}

func convertGetMeResponse(user *domain.User) *GetMeResponse {
	return &GetMeResponse{
		Me: convertUserResponse(user),
	}
}
