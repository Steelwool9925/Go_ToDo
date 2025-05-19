package cmd

import (
	pb "Go_Test/api"
	"Go_Test/client"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// getTasksCmd represents the command to fetch and display all tasks.
var getTasksCmd = &cobra.Command{
	Use:   "get-tasks",
	Short: "Fetches and displays the list of all tasks from the server",
	Long:  `Connects to the gRPC server, calls the GetTasks RPC method, and prints the results.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			commonFxOptions(),
			client.Module,
			fx.Invoke(runGetTasksLogic),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := app.Start(ctx); err != nil {
			return fmt.Errorf("fx app failed to start: %w", err)
		}
		if err := app.Stop(ctx); err != nil {
			return fmt.Errorf("fx app failed to stop gracefully: %w", err)
		}
		return nil
	},
}

func runGetTasksLogic(lc fx.Lifecycle, taskClient pb.TaskServiceClient, logger *zap.Logger) {
	logger.Info("Executing GetTasks logic via CLI command")
	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasksReply, err := taskClient.GetTasks(reqCtx, &pb.GetTasksRequest{})
	if err != nil {
		logger.Error("Failed to get tasks via CLI", zap.Error(err))
		fmt.Printf("Error: Could not get tasks: %v\n", err)
		return
	}
	logger.Info("Tasks received successfully via CLI", zap.Int("count", len(tasksReply.GetTasks())))
	if len(tasksReply.GetTasks()) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Println("--- Tasks ---")
	for i, task := range tasksReply.GetTasks() {
		fmt.Printf("%d. ID: %s\n", i+1, task.GetId())
		fmt.Printf("   Title: %s\n", task.GetTitle())
		fmt.Printf("   Description: %s\n", task.GetDescription())
		fmt.Printf("   Status: %s\n", task.GetStatus())
		fmt.Printf("   Created At: %s\n", task.GetCreatedAt())
		fmt.Printf("   Updated At: %s\n", task.GetUpdatedAt())
		fmt.Println("---------------")
	}
}

func init() {
	clientCmd.AddCommand(getTasksCmd)
}
