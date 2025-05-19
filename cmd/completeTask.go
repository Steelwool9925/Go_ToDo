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

var (
	completeTaskID string
)

// completeTaskCmd represents the command to mark a task as completed.
var completeTaskCmd = &cobra.Command{
	Use:   "complete-task --id <task_id>",
	Short: "Marks a specified task as completed",
	Long:  `Connects to the gRPC server and calls the CompleteTask RPC method for the given task ID. Handles errors for non-existent or already completed tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if completeTaskID == "" {
			return fmt.Errorf("task ID is required. Use --id flag")
		}

		app := fx.New(
			commonFxOptions(),
			client.Module,
			fx.Supply(
				&pb.CompleteTaskRequest{
					TaskId: completeTaskID,
				},
			),
			fx.Invoke(runCompleteTaskLogic),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Start(ctx); err != nil {
			return fmt.Errorf("fx app failed to start for complete-task: %w", err)
		}
		if err := app.Stop(ctx); err != nil {
			return fmt.Errorf("fx app failed to stop gracefully for complete-task: %w", err)
		}
		return nil
	},
}

func runCompleteTaskLogic(lc fx.Lifecycle, taskClient pb.TaskServiceClient, logger *zap.Logger, req *pb.CompleteTaskRequest) {
	logger.Info("Executing CompleteTask logic via CLI command",
		zap.String("task_id", req.GetTaskId()))

	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reply, err := taskClient.CompleteTask(reqCtx, req)
	if err != nil {
		logger.Error("Failed to complete task via CLI", zap.Error(err))
		fmt.Printf("Error completing task: %v\n", err)
		return
	}

	completedTask := reply.GetTask()
	logger.Info("Task completed successfully via CLI", zap.String("id", completedTask.GetId()))
	fmt.Println("--- Task Completed Successfully ---")
	fmt.Printf("ID: %s\n", completedTask.GetId())
	fmt.Printf("Title: %s\n", completedTask.GetTitle())
	fmt.Printf("Status: %s\n", completedTask.GetStatus())
	fmt.Printf("Updated At: %s\n", completedTask.GetUpdatedAt())
	fmt.Println("-------------------------------")
}

func init() {
	completeTaskCmd.Flags().StringVar(&completeTaskID, "id", "", "ID of the task to complete (required)")
	clientCmd.AddCommand(completeTaskCmd)
}
