package server

import (
	pb "Go_Test/api"
	cfg "Go_Test/config"
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Module exports providers for the gRPC server and its TaskService implementation for FX.
var Module = fx.Options(
	fx.Provide(NewGRPCServer),
	fx.Provide(NewTaskServiceImpl),
)

type GRPCServerParams struct {
	fx.In
	Lifecycle         fx.Lifecycle
	Logger            *zap.Logger
	Config            *cfg.Config
	TaskServiceServer pb.TaskServiceServer
}

// NewGRPCServer creates, configures, and manages the lifecycle of the main gRPC server.
func NewGRPCServer(p GRPCServerParams) (*grpc.Server, error) {
	p.Logger.Info("Setting up gRPC server for TaskService")

	var serverOpts []grpc.ServerOption
	server := grpc.NewServer(serverOpts...)

	pb.RegisterTaskServiceServer(server, p.TaskServiceServer)
	reflection.Register(server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pb.TaskService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("Starting gRPC server", zap.String("address", p.Config.GRPCServerAddress))
			lis, err := net.Listen("tcp", p.Config.GRPCServerAddress)
			if err != nil {
				p.Logger.Error("Failed to listen for gRPC", zap.Error(err))
				return fmt.Errorf("failed to listen for gRPC: %w", err)
			}
			go func() {
				if err := server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
					p.Logger.Error("gRPC server failed to serve", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping gRPC server")
			server.GracefulStop()
			p.Logger.Info("gRPC server stopped")
			return nil
		},
	})

	return server, nil
}
