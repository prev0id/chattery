package utils

import (
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chattery/backend/user_service/internal/domain"
)

func HandleGRPCError(err error) error {
	if err == nil {
		return nil
	}
	if code, found := unwrapCode(err); found {
		return status.Error(code, err.Error())
	}
	slog.Error("[error_util] unknown error", slog.String("error", err.Error()))
	return err
}

func unwrapCode(err error) (codes.Code, bool) {
	for targetErr, code := range domain.ErrToGRPCCode {
		if errors.Is(err, targetErr) {
			return code, true
		}
	}
	return codes.Internal, false
}
