package user_api

import (
	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type PostCreateUserRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PostCreateUserResponse struct {
	ID int64 `json:"id"`
}

type PostLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUsersResponse struct {
	Users []User
}

type PostUpdateUserRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetMeResponse struct {
	Me User `json:"me"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func convertPostCreateUserRequest(req *PostCreateUserRequest) *domain.User {
	login := domain.Login(req.Login)
	return &domain.User{
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}

func convertPostCreateUserResponse(userID domain.UserID) *PostCreateUserResponse {
	return &PostCreateUserResponse{
		ID: userID.I64(),
	}
}

func convertGetUsersResponse(users []*domain.User) GetUsersResponse {
	return GetUsersResponse{
		Users: sliceutil.Map(users, convertUserResponse),
	}
}

func convertPostUpdateUserRequest(req *PostUpdateUserRequest, userID domain.UserID) *domain.User {
	login := domain.Login(req.Login)
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
		Avatar:   "/v1/image/" + user.Username.String() + ".png",
		// AvatarURL: "/image/" + user.AvatarID.String(),
	}
}

func convertGetMeResponse(user *domain.User) *GetMeResponse {
	return &GetMeResponse{
		Me: convertUserResponse(user),
	}
}
