package cmd

import (
	"Go_Test/database"
	"Go_Test/repository"
	"Go_Test/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// serverCmd represents the command to start the gRPC server.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the gRPC TaskService server",
	Long:  `Initializes and runs the gRPC server, including database connections and service implementations, using Uber FX.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			commonFxOptions(),
			database.Module,
			repository.Module,
			server.Module,
			fx.Invoke(func(*grpc.Server, *zap.Logger) {}), // Ensure server and logger are initialized
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := app.Start(ctx); err != nil {
			return err
		}

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
		<-stopChan

		if err := app.Stop(context.Background()); err != nil {
			return err
		}
		return nil
	},
}
