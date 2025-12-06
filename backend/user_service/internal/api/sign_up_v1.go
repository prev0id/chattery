package api

import (
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
	"context"
)

func (s *Server) SignUpV1(ctx context.Context, request *user_servicepb.SignUpV1Request) (*user_servicepb.SignUpV1Response, error) {
	return nil, nil
}
