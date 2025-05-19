package server

import (
	pb "Go_Test/api"
	repo "Go_Test/repository"
	"context"
	"database/sql"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskServiceImpl implements the proto.TaskServiceServer interface for task-related RPC calls.
type TaskServiceImpl struct {
	pb.UnimplementedTaskServiceServer
	logger   *zap.Logger
	taskRepo repo.TaskRepository
}

// NewTaskServiceImpl creates a new TaskServiceImpl.
func NewTaskServiceImpl(logger *zap.Logger, taskRepo repo.TaskRepository) pb.TaskServiceServer {
	return &TaskServiceImpl{logger: logger, taskRepo: taskRepo}
}

// GetTasks handles the RPC call to fetch all tasks.
func (s *TaskServiceImpl) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksReply, error) {
	s.logger.Info("TaskServiceImpl: GetTasks called")
	tasks, err := s.taskRepo.FetchTasks(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch tasks in service", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to fetch tasks: %v", err)
	}
	return &pb.GetTasksReply{Tasks: tasks}, nil
}

// AddTask handles the RPC call to add a new task.
func (s *TaskServiceImpl) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskReply, error) {
	s.logger.Info("TaskServiceImpl: AddTask called", zap.String("title", req.GetTitle()))
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title cannot be empty")
	}
	taskStatus := req.GetStatus()
	if taskStatus == "" {
		taskStatus = "pending"
	}
	createdTask, err := s.taskRepo.AddTask(ctx, req.GetTitle(), req.GetDescription(), taskStatus)
	if err != nil {
		s.logger.Error("Failed to add task in service", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add task: %v", err)
	}
	return &pb.AddTaskReply{Task: createdTask}, nil
}

// CompleteTask handles the RPC call to mark a task as completed.
// It includes error handling for non-existent tasks or tasks already completed.
func (s *TaskServiceImpl) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskReply, error) {
	s.logger.Info("TaskServiceImpl: CompleteTask called", zap.String("task_id", req.GetTaskId()))
	if req.GetTaskId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "task_id cannot be empty")
	}

	existingTask, err := s.taskRepo.FetchTaskByID(ctx, req.GetTaskId())
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("CompleteTask: Task not found", zap.String("task_id", req.GetTaskId()))
			return nil, status.Errorf(codes.NotFound, "task with ID '%s' not found", req.GetTaskId())
		}
		s.logger.Error("CompleteTask: Failed to fetch task", zap.String("task_id", req.GetTaskId()), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve task details: %v", err)
	}

	if existingTask.GetStatus() == "completed" {
		s.logger.Info("CompleteTask: Task already completed", zap.String("task_id", req.GetTaskId()))
		return nil, status.Errorf(codes.FailedPrecondition, "task with ID '%s' is already completed", req.GetTaskId())
	}

	updatedTask, err := s.taskRepo.UpdateTaskStatus(ctx, req.GetTaskId(), "completed")
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn("CompleteTask: Task disappeared before update", zap.String("task_id", req.GetTaskId()))
			return nil, status.Errorf(codes.NotFound, "task with ID '%s' not found for update", req.GetTaskId())
		}
		s.logger.Error("CompleteTask: Failed to update task status", zap.String("task_id", req.GetTaskId()), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to complete task: %v", err)
	}

	s.logger.Info("TaskServiceImpl: Task completed successfully", zap.String("task_id", updatedTask.GetId()))
	return &pb.CompleteTaskReply{Task: updatedTask}, nil
}
