package api

import (
	"context"

	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"chattery/backend/user_service/internal/utils"
)

func (s *Server) UpdateV1(ctx context.Context, request *user_servicepb.UpdateV1Request) (*user_servicepb.UpdateV1Response, error) {
	user := convertUserFromUpdateRequest(request)
	if err := s.service.Update(ctx, request.GetSession().GetValue(), user); err != nil {
		return nil, utils.HandleGRPCError(err)
	}
	return &user_servicepb.UpdateV1Response{}, nil
}
