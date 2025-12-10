package api

import (
	"context"

	user_servicepb "chattery/backend/common/pb/user_service"
	"chattery/backend/user_service/internal/utils"
)

func (s *Server) LogInV1(ctx context.Context, request *user_servicepb.LogInV1Request) (*user_servicepb.LogInV1Response, error) {
	session, err := s.service.Login(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, utils.HandleGRPCError(err)
	}
	return &user_servicepb.LogInV1Response{
		Session: convertSessionToPB(session),
	}, nil
}
