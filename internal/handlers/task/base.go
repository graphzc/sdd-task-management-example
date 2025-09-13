package task

import (
	"context"

	"github.com/graphzc/sdd-task-management-example/internal/dto"
	"github.com/graphzc/sdd-task-management-example/internal/services/task"
	"github.com/graphzc/sdd-task-management-example/internal/utils/echoutil"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
)

type Handler interface {
	CreateTask(ctx context.Context, req *dto.TaskCreateRequest, userID string) (*dto.MessageResponse, error)
	GetTaskByID(ctx context.Context, taskID string, userID string) (*dto.TaskResponse, error)
	GetTasksByUserID(ctx context.Context, userID string) ([]dto.TaskResponse, error)
	UpdateTaskByID(ctx context.Context, taskID string, req *dto.TaskUpdateRequest, userID string) (*dto.MessageResponse, error)
	UpdateTaskStatusByID(ctx context.Context, taskID string, req *dto.TaskUpdateStatusRequest, userID string) (*dto.MessageResponse, error)
	DeleteTaskByID(ctx context.Context, taskID string, userID string) (*dto.MessageResponse, error)

	// Wrapper methods for WrapWithStatus compatibility
	CreateTaskWrapped(ctx context.Context, req *dto.TaskCreateRequest) (*dto.MessageResponse, error)
	GetTaskByIDWrapped(ctx context.Context, req *dto.TaskGetByIDRequest) (*dto.TaskResponse, error)
	GetTasksByUserIDWrapped(ctx context.Context, _ any) ([]dto.TaskResponse, error)
	UpdateTaskByIDWrapped(ctx context.Context, req *dto.TaskUpdateWithIDRequest) (*dto.MessageResponse, error)
	UpdateTaskStatusByIDWrapped(ctx context.Context, req *dto.TaskUpdateStatusWithIDRequest) (*dto.MessageResponse, error)
	DeleteTaskByIDWrapped(ctx context.Context, req *dto.TaskDeleteRequest) (*dto.MessageResponse, error)
}

type handler struct {
	taskService task.Service
}

// @WireSet("Handler")
func New(taskService task.Service) Handler {
	return &handler{
		taskService: taskService,
	}
}

func (h *handler) CreateTask(ctx context.Context, req *dto.TaskCreateRequest, userID string) (*dto.MessageResponse, error) {
	serviceInput := task.TaskCreateInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
	}

	err := h.taskService.CreateTask(ctx, &serviceInput, userID)
	if err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Task created successfully",
	}, nil
}

func (h *handler) GetTaskByID(ctx context.Context, taskID string, userID string) (*dto.TaskResponse, error) {
	foundTask, err := h.taskService.FindTaskByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.TaskResponse{
		ID:          foundTask.ID,
		UserID:      foundTask.UserID,
		Title:       foundTask.Title,
		Description: foundTask.Description,
		Priority:    foundTask.Priority,
		Status:      foundTask.Status,
		CreatedAt:   foundTask.CreatedAt,
		UpdatedAt:   foundTask.UpdatedAt,
	}, nil
}

func (h *handler) GetTasksByUserID(ctx context.Context, userID string) ([]dto.TaskResponse, error) {
	tasks, err := h.taskService.FindTaskByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = dto.TaskResponse{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return taskResponses, nil
}

func (h *handler) UpdateTaskByID(ctx context.Context, taskID string, req *dto.TaskUpdateRequest, userID string) (*dto.MessageResponse, error) {
	serviceInput := task.TaskUpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
	}

	err := h.taskService.UpdateTaskByID(ctx, taskID, &serviceInput, userID)
	if err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Task updated successfully",
	}, nil
}

func (h *handler) UpdateTaskStatusByID(ctx context.Context, taskID string, req *dto.TaskUpdateStatusRequest, userID string) (*dto.MessageResponse, error) {
	serviceInput := task.TaskUpdateStatusInput{
		Status: req.Status,
	}

	err := h.taskService.UpdateTaskStatusByID(ctx, taskID, &serviceInput, userID)
	if err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Task status updated successfully",
	}, nil
}

func (h *handler) DeleteTaskByID(ctx context.Context, taskID string, userID string) (*dto.MessageResponse, error) {
	err := h.taskService.DeleteTaskByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Task deleted successfully",
	}, nil
}

// Wrapper methods for WrapWithStatus compatibility

func (h *handler) CreateTaskWrapped(ctx context.Context, req *dto.TaskCreateRequest) (*dto.MessageResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	return h.CreateTask(ctx, req, userID)
}

func (h *handler) GetTaskByIDWrapped(ctx context.Context, req *dto.TaskGetByIDRequest) (*dto.TaskResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	return h.GetTaskByID(ctx, req.ID, userID)
}

func (h *handler) GetTasksByUserIDWrapped(ctx context.Context, _ any) ([]dto.TaskResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	return h.GetTasksByUserID(ctx, userID)
}

func (h *handler) UpdateTaskByIDWrapped(ctx context.Context, req *dto.TaskUpdateWithIDRequest) (*dto.MessageResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	updateReq := &dto.TaskUpdateRequest{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
	}
	return h.UpdateTaskByID(ctx, req.ID, updateReq, userID)
}

func (h *handler) UpdateTaskStatusByIDWrapped(ctx context.Context, req *dto.TaskUpdateStatusWithIDRequest) (*dto.MessageResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	statusReq := &dto.TaskUpdateStatusRequest{
		Status: req.Status,
	}
	return h.UpdateTaskStatusByID(ctx, req.ID, statusReq, userID)
}

func (h *handler) DeleteTaskByIDWrapped(ctx context.Context, req *dto.TaskDeleteRequest) (*dto.MessageResponse, error) {
	userID, err := echoutil.GetUserIDFromContext(ctx)
	if err != nil {
return nil, servererr.NewError(
			servererr.ErrorCodeUnauthorized,
			"user ID not found in context",
		)
	}
	return h.DeleteTaskByID(ctx, req.ID, userID)
}
