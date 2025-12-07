package api

import (
	"context"

	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"chattery/backend/user_service/internal/utils"
)

func (s *Server) ValidateV1(ctx context.Context, request *user_servicepb.ValidateV1Request) (*user_servicepb.ValidateV1Response, error) {
	claims, err := s.service.Validate(request.GetSession().GetValue())
	if err != nil {
		return nil, utils.HandleGRPCError(err)
	}
	return &user_servicepb.ValidateV1Response{
		Info: convertClaimsToUserInfo(claims),
	}, nil
}
