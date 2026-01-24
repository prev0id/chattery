package user_api

import (
	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

type CreateRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SearchResponse struct {
	Users []User
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func convertCreateRequest(req *CreateRequest) *domain.User {
	login := domain.Login(req.Login)
	return &domain.User{
		Username: domain.Username(req.Username),
		Login:    login,
		Password: domain.NewPassword(req.Password, login),
	}
}

func converSearchResponse(users []*domain.User) SearchResponse {
	return SearchResponse{
		Users: sliceutil.Map(users, convertUserResponse),
	}

}

func convertUserResponse(user *domain.User) User {
	return User{
		ID:        user.ID.I64(),
		Username:  user.Username.String(),
		AvatarURL: "/image/" + user.AvatarID.String(),
	}
}

// type UserCursor struct {
// 	ID        int64     `json:"id"`
// 	Timestamp time.Time `json:"timestamp"`
// }

// func convertCursorRequest(req UserCursor) *domain.UserCursor {
// 	return &domain.UserCursor{
// 		ID:        domain.UserID(req.ID),
// 		Timestamp: req.Timestamp,
// 	}
// }

// func converCursorResponse(cursor *domain.UserCursor) UserCursor {
// 	return UserCursor{
// 		ID:        cursor.ID.I64(),
// 		Timestamp: cursor.Timestamp,
// 	}
// }
