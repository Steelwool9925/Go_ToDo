package repository

import (
	pb "Go_Test/api"
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module exports the TaskRepository provider for FX.
var Module = fx.Options(
	fx.Provide(NewSQLTaskRepository),
)

// TaskRepository defines the interface for task data persistence operations.
type TaskRepository interface {
	FetchTasks(ctx context.Context) ([]*pb.Task, error)
	AddTask(ctx context.Context, title string, description string, status string) (*pb.Task, error)
	FetchTaskByID(ctx context.Context, taskID string) (*pb.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, newStatus string) (*pb.Task, error)
}

type sqlTaskRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewSQLTaskRepository creates a new SQL-based task repository.
func NewSQLTaskRepository(db *sql.DB, logger *zap.Logger) TaskRepository {
	return &sqlTaskRepository{db: db, logger: logger}
}

// FetchTasks retrieves all tasks from the database.
func (r *sqlTaskRepository) FetchTasks(ctx context.Context) ([]*pb.Task, error) {
	r.logger.Debug("Fetching tasks from database")
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("Failed to query tasks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var tasks []*pb.Task
	for rows.Next() {
		var task pb.Task
		var createdAt, updatedAt sql.NullTime
		var description sql.NullString
		if err := rows.Scan(&task.Id, &task.Title, &description, &task.Status, &createdAt, &updatedAt); err != nil {
			r.logger.Error("Failed to scan task row", zap.Error(err))
			return nil, err
		}
		if description.Valid {
			task.Description = description.String
		} else {
			task.Description = ""
		}
		if createdAt.Valid {
			task.CreatedAt = createdAt.Time.Format(time.RFC3339)
		}
		if updatedAt.Valid {
			task.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		r.logger.Error("Error during rows iteration for tasks", zap.Error(err))
		return nil, err
	}
	r.logger.Debug("Successfully fetched tasks", zap.Int("count", len(tasks)))
	return tasks, nil
}

// AddTask inserts a new task into the database and returns the created task.
func (r *sqlTaskRepository) AddTask(ctx context.Context, title string, description string, status string) (*pb.Task, error) {
	r.logger.Debug("Adding new task to database", zap.String("title", title))
	query := "INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, title, sql.NullString{String: description, Valid: description != ""}, status)
	if err != nil {
		r.logger.Error("Failed to insert task", zap.Error(err))
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("Failed to get last insert ID for task", zap.Error(err))
		return nil, err
	}
	return r.FetchTaskByID(ctx, fmt.Sprintf("%d", id))
}

// FetchTaskByID retrieves a single task by its ID.
func (r *sqlTaskRepository) FetchTaskByID(ctx context.Context, taskID string) (*pb.Task, error) {
	r.logger.Debug("Fetching task by ID", zap.String("taskID", taskID))
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?"
	var task pb.Task
	var createdAt, updatedAt sql.NullTime
	var description sql.NullString
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&task.Id, &task.Title, &description, &task.Status, &createdAt, &updatedAt,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error("Failed to fetch task by ID", zap.String("taskID", taskID), zap.Error(err))
		}
		return nil, err
	}
	if description.Valid {
		task.Description = description.String
	} else {
		task.Description = ""
	}
	if createdAt.Valid {
		task.CreatedAt = createdAt.Time.Format(time.RFC3339)
	}
	if updatedAt.Valid {
		task.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}
	r.logger.Debug("Successfully fetched task by ID", zap.String("taskID", taskID))
	return &task, nil
}

// UpdateTaskStatus updates the status of a task and returns the updated task.
func (r *sqlTaskRepository) UpdateTaskStatus(ctx context.Context, taskID string, newStatus string) (*pb.Task, error) {
	r.logger.Debug("Updating task status", zap.String("taskID", taskID), zap.String("newStatus", newStatus))
	query := "UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, newStatus, taskID)
	if err != nil {
		r.logger.Error("Failed to update task status", zap.String("taskID", taskID), zap.Error(err))
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected after status update", zap.String("taskID", taskID), zap.Error(err))
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	return r.FetchTaskByID(ctx, taskID)
}
