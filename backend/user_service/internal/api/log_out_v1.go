package api

import (
	"context"

	user_servicepb "chattery/backend/common/pb/user_service"
)

func (s *Server) LogOutV1(ctx context.Context, request *user_servicepb.LogOutV1Request) (*user_servicepb.LogOutV1Response, error) {
	session := s.service.CreateSessionInvalidation()
	return &user_servicepb.LogOutV1Response{
		Session: convertSessionToPB(session),
	}, nil
}
