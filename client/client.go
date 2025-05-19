package client

import (
	pb "Go_Test/api"
	cfg "Go_Test/config"
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Module exports providers for the gRPC client connection and TaskService client for FX.
var Module = fx.Options(
	fx.Provide(NewGRPCConnection),
	fx.Provide(NewTaskServiceClient),
)

type GRPCConnectionParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Config    *cfg.Config
}

// NewGRPCConnection creates and manages the lifecycle of a gRPC client connection.
func NewGRPCConnection(p GRPCConnectionParams) (*grpc.ClientConn, error) {
	p.Logger.Info("Setting up gRPC client connection", zap.String("target", p.Config.GRPCClientTarget))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	conn, err := grpc.DialContext(context.Background(), p.Config.GRPCClientTarget, opts...)
	if err != nil {
		p.Logger.Error("Failed to dial gRPC server", zap.Error(err))
		return nil, fmt.Errorf("failed to dial gRPC server %s: %w", p.Config.GRPCClientTarget, err)
	}
	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Closing gRPC client connection")
			if err := conn.Close(); err != nil {
				p.Logger.Error("Failed to close gRPC client connection", zap.Error(err))
				return err
			}
			p.Logger.Info("gRPC client connection closed")
			return nil
		},
	})
	return conn, nil
}

// NewTaskServiceClient creates a new TaskService client stub.
func NewTaskServiceClient(conn *grpc.ClientConn) pb.TaskServiceClient {
	return pb.NewTaskServiceClient(conn)
}
