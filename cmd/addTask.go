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
	taskTitle       string
	taskDescription string
	taskStatus      string
)

// addTaskCmd represents the command to add a new task.
var addTaskCmd = &cobra.Command{
	Use:   "add-task --title <title> [--description <desc>] [--status <status>]",
	Short: "Adds a new task via the gRPC server",
	Long:  `Connects to the gRPC server and calls the AddTask RPC method with the provided details to create a new task.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if taskTitle == "" {
			return fmt.Errorf("title is required. Use --title or -t flag")
		}

		app := fx.New(
			commonFxOptions(),
			client.Module,
			fx.Supply(
				&pb.AddTaskRequest{
					Title:       taskTitle,
					Description: taskDescription,
					Status:      taskStatus,
				},
			),
			fx.Invoke(runAddTaskLogic),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := app.Start(ctx); err != nil {
			return fmt.Errorf("fx app failed to start for add-task: %w", err)
		}
		if err := app.Stop(ctx); err != nil {
			return fmt.Errorf("fx app failed to stop gracefully for add-task: %w", err)
		}
		return nil
	},
}

func runAddTaskLogic(lc fx.Lifecycle, taskClient pb.TaskServiceClient, logger *zap.Logger, req *pb.AddTaskRequest) {
	logger.Info("Executing AddTask logic via CLI command",
		zap.String("title", req.GetTitle()),
		zap.String("description", req.GetDescription()),
		zap.String("status", req.GetStatus()))

	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reply, err := taskClient.AddTask(reqCtx, req)
	if err != nil {
		logger.Error("Failed to add task via CLI", zap.Error(err))
		fmt.Printf("Error: Could not add task: %v\n", err)
		return
	}

	createdTask := reply.GetTask()
	logger.Info("Task added successfully via CLI", zap.String("id", createdTask.GetId()))
	fmt.Println("--- Task Added Successfully ---")
	fmt.Printf("ID: %s\n", createdTask.GetId())
	fmt.Printf("Title: %s\n", createdTask.GetTitle())
	fmt.Printf("Description: %s\n", createdTask.GetDescription())
	fmt.Printf("Status: %s\n", createdTask.GetStatus())
	fmt.Printf("Created At: %s\n", createdTask.GetCreatedAt())
	fmt.Printf("Updated At: %s\n", createdTask.GetUpdatedAt())
	fmt.Println("-----------------------------")
}

func init() {
	addTaskCmd.Flags().StringVarP(&taskTitle, "title", "t", "", "Title of the task (required)")
	addTaskCmd.Flags().StringVarP(&taskDescription, "description", "d", "", "Description of the task")
	addTaskCmd.Flags().StringVarP(&taskStatus, "status", "s", "", "Status of the task (e.g., pending, in_progress). Defaults to 'pending' server-side if empty.")
	clientCmd.AddCommand(addTaskCmd)
}
