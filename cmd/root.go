package cmd

import (
	"Go_Test/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands.
// It serves as the entry point for the CLI application.
var rootCmd = &cobra.Command{
	Use:   "fx-grpc-app",
	Short: "A gRPC task management application with Uber FX and Cobra.",
	Long:  `This application provides a gRPC server and client CLI for managing tasks. Built with Go, Uber FX, and Cobra.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// commonFxOptions provides shared FX options for logger and configuration,
// used by multiple CLI commands.
func commonFxOptions() fx.Option {
	return fx.Options(
		fx.Provide(NewLogger),
		config.Module,
	)
}

// NewLogger provides a zap logger instance for FX.
func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(clientCmd)
}
