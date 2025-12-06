package api

import (
	"chattery/backend/user_service/internal/domain"
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
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
