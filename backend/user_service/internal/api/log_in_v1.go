package api

import (
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"context"
)

func (s *Server) LogInV1(ctx context.Context, request *user_servicepb.LogInV1Request) (*user_servicepb.LogInV1Response, error) {
	return nil, nil
}
