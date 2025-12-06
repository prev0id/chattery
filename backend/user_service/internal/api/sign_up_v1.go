package api

import (
	"context"

	"chattery/backend/user_service/internal/domain"
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
)

func (s *Server) SignUpV1(ctx context.Context, request *user_servicepb.SignUpV1Request) (*user_servicepb.SignUpV1Response, error) {
	user := convertUserFromSignUpRequest(request)

	session, err := s.service.CreateUser(ctx, user)
	if err != nil {
		return nil, domain.HandleGRPCError(err)
	}

	return &user_servicepb.SignUpV1Response{
		Session: convertSessionToPB(session),
	}, nil
}
