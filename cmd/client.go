package cmd

import (
	"github.com/spf13/cobra"
)

// clientCmd represents the base command for client-side operations.
// It groups various client actions like get-tasks, add-task, complete-task.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Manages client-side gRPC operations for TaskService",
	Long:  `A parent command for various client actions interacting with the TaskService.`,
}

func init() {
}
