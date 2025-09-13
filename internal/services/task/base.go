package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/graphzc/sdd-task-management-example/internal/domain/entities"
	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
	"github.com/graphzc/sdd-task-management-example/internal/repositories/task"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/graphzc/sdd-task-management-example/internal/utils/timeutil"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateTask(ctx context.Context, in *TaskCreateInput, userID string) error
	FindTaskByID(ctx context.Context, taskID string, userID string) (*entities.Task, error)
	FindTaskByUserID(ctx context.Context, userID string) ([]entities.Task, error)
	DeleteTaskByID(ctx context.Context, taskID string, userID string) error
	UpdateTaskByID(ctx context.Context, taskID string, in *TaskUpdateInput, userID string) error
	UpdateTaskStatusByID(ctx context.Context, taskID string, in *TaskUpdateStatusInput, userID string) error
}

type service struct {
	config   *config.Config
	taskRepo task.Repository
}

// @WireSet("Service")
func NewService(
	config *config.Config,
	taskRepo task.Repository,
) Service {
	return &service{
		config:   config,
		taskRepo: taskRepo,
	}
}

func (s *service) CreateTask(ctx context.Context, in *TaskCreateInput, userID string) error {
	// Create new task entity
	newTask := &entities.Task{
		ID:          uuid.NewString(),
		UserID:      userID,
		Title:       in.Title,
		Description: in.Description,
		Priority:    enums.TaskPriority(in.Priority),
		Status:      enums.TaskStatusTodo, // Default status
		CreatedAt:   timeutil.BangkokNow(),
		UpdatedAt:   timeutil.BangkokNow(),
	}

	// Create task in repository
	_, err := s.taskRepo.Create(ctx, newTask)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to create task")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to create task",
		)
	}

	return nil
}

func (s *service) FindTaskByID(ctx context.Context, taskID string, userID string) (*entities.Task, error) {
	// Find the task
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		log.Error().
			Err(err).
			Str("taskId", taskID).
			Msg("Failed to find task by ID")

		return nil, servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to find task",
		)
	}

	// Check if task exists
	if task == nil {
		log.Warn().
			Str("taskId", taskID).
			Msg("Task not found")

		return nil, servererr.NewError(
			servererr.ErrorCodeNotFound,
			"Task not found",
		)
	}

	// Check if task belongs to the user
	if task.UserID != userID {
		log.Warn().
			Str("taskId", taskID).
			Str("userId", userID).
			Str("taskUserId", task.UserID).
			Msg("Task does not belong to user")

		return nil, servererr.NewError(
			servererr.ErrorCodeNotFound,
			"Task not found",
		)
	}

	return task, nil
}

func (s *service) FindTaskByUserID(ctx context.Context, userID string) ([]entities.Task, error) {
	tasks, err := s.taskRepo.FindByUserID(ctx, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("userId", userID).
			Msg("Failed to find tasks by user ID")

		return nil, servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to find tasks",
		)
	}

	return tasks, nil
}

func (s *service) DeleteTaskByID(ctx context.Context, taskID string, userID string) error {
	// Find the task first
	_, err := s.FindTaskByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	// Delete the task
	if err := s.taskRepo.DeleteByID(ctx, taskID); err != nil {
		log.Error().
			Err(err).
			Str("taskId", taskID).
			Msg("Failed to delete task")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to delete task",
		)
	}

	return nil
}

func (s *service) UpdateTaskByID(ctx context.Context, taskID string, in *TaskUpdateInput, userID string) error {
	// Validate priority
	if in.Priority < 1 || in.Priority > 3 {
		log.Warn().
			Int("priority", in.Priority).
			Msg("Invalid task priority")

		return servererr.NewError(
			servererr.ErrorCodeBadRequest,
			"Invalid priority. Priority must be between 1 and 3",
		)
	}

	// Find the task first to ensure it exists and belongs to user
	_, err := s.FindTaskByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	// Update task in repository
	if err := s.taskRepo.UpdateByID(ctx, taskID, in.Title, in.Description, enums.TaskPriority(in.Priority)); err != nil {
		log.Error().
			Err(err).
			Str("taskId", taskID).
			Msg("Failed to update task")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to update task",
		)
	}

	return nil
}

func (s *service) UpdateTaskStatusByID(ctx context.Context, taskID string, in *TaskUpdateStatusInput, userID string) error {
	// Validate status
	statusEnum := enums.TaskStatus(in.Status)
	if statusEnum != enums.TaskStatusTodo && statusEnum != enums.TaskStatusInProgress && statusEnum != enums.TaskStatusCompleted {
		log.Warn().
			Str("status", in.Status).
			Msg("Invalid task status")

		return servererr.NewError(
			servererr.ErrorCodeBadRequest,
			"Invalid status. Status must be TODO, IN_PROGRESS, or COMPLETED",
		)
	}

	// Find the task first to ensure it exists and belongs to user
	_, err := s.FindTaskByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	// Update task status in repository
	if err := s.taskRepo.UpdateStatusByID(ctx, taskID, statusEnum); err != nil {
		log.Error().
			Err(err).
			Str("taskId", taskID).
			Msg("Failed to update task status")

		return servererr.NewError(
			servererr.ErrorCodeInternalServerError,
			"Failed to update task status",
		)
	}

	return nil
}
