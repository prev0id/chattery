package api

import (
	user_servicepb "chattery/backend/common/pb/user_service"
	"chattery/backend/user_service/internal/domain"
)

func convertUserFromSignUpRequest(request *user_servicepb.SignUpV1Request) domain.User {
	user := domain.User{
		Username: request.GetUsername(),
		Login:    request.GetLogin(),
	}

	user.SetRawPassword(request.GetPassword())

	return user
}

func convertSessionToPB(session domain.Session) *user_servicepb.Session {
	return &user_servicepb.Session{
		Name:     session.Name,
		Value:    session.Value,
		Path:     session.Path,
		MaxAge:   int64(session.MaxAge),
		Secure:   session.Secure,
		HttpOnly: session.HttpOnly,
		SameSite: int64(session.SameSite),
	}
}

func convertClaimsToUserInfo(claims domain.Claims) *user_servicepb.UserInfo {
	return &user_servicepb.UserInfo{
		UserId:   int64(claims.UserID),
		Username: claims.Username,
		ImageId:  string(claims.ImageID),
	}
}

func convertUserFromUpdateRequest(request *user_servicepb.UpdateV1Request) domain.User {
	user := domain.User{
		Login:    request.GetLogin(),
		Username: request.GetUsername(),
	}
	if request.GetLogin() != "" && request.GetPassword() != "" {
		user.SetRawPassword(request.GetPassword())
	}
	return user
}
