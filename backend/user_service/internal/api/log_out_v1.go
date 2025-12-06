package api

import (
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"context"
)

func (s *Server) LogOutV1(ctx context.Context, request *user_servicepb.LogOutV1Request) (*user_servicepb.LogOutV1Response, error) {
	return nil, nil
}
