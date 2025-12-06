package api

import (
	"chattery/backend/user_service/internal/domain"
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"context"
)

func (s *Server) LogInV1(ctx context.Context, request *user_servicepb.LogInV1Request) (*user_servicepb.LogInV1Response, error) {
	session, err := s.service.Login(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, domain.HandleGRPCError(err)
	}
	return &user_servicepb.LogInV1Response{
		Session: convertSessionToPB(session),
	}, nil
}
