package api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"chattery/backend/user_service/internal/config"
	"chattery/backend/user_service/internal/pb"
	user_servicepb "chattery/backend/user_service/internal/pb/user_service"
)

type Server struct {
	user_servicepb.UnimplementedUserServiceServer
}

func StartChatServer(ctx context.Context) error {
	server := &Server{}

	if err := server.handleGRPC(); err != nil {
		return fmt.Errorf("server.handleGRPC: %w", err)
	}

	if err := server.handleHTTP(ctx); err != nil {
		return fmt.Errorf("server.handleHTTP: %w", err)
	}

	return nil
}

func (s *Server) handleGRPC() error {
	grpcServer := grpc.NewServer()
	user_servicepb.RegisterUserServiceServer(grpcServer, s)

	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		return fmt.Errorf("net.Listen address=%s: %w", config.GRPCAddress, err)
	}

	go func() {
		slog.Info("starting GRPC server", slog.String("address", config.GRPCAddress))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("grpcServer.Serve: %s", err.Error())
		}
	}()

	return nil
}

func (s *Server) handleHTTP(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := user_servicepb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.GRPCAddress, opts); err != nil {
		return fmt.Errorf("chat_servicepb.RegisterChatServiceHandlerFromEndpoint: %w", err)
	}

	pb.RegisterSwaggerHandlers(mux)

	go func() {
		slog.Info("starting HTTP server", slog.String("address", config.HTTPAddress))
		if err := http.ListenAndServe(config.HTTPAddress, mux); err != nil {
			log.Fatalf("http.ListenAndServe: %s", err.Error())
		}
	}()

	return nil
}
