package dto

import (
	"time"

	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
)

type TaskCreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    int    `json:"priority" validate:"required,min=1,max=3"`
}

type TaskUpdateRequest = TaskCreateRequest

type TaskUpdateStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type TaskResponse struct {
	ID          string             `json:"id"`
	UserID      string             `json:"userId"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Priority    enums.TaskPriority `json:"priority"`
	Status      enums.TaskStatus   `json:"status"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// Request DTOs for wrapped handlers
type TaskGetByIDRequest struct {
	ID string `param:"id" validate:"required"`
}

type TaskUpdateWithIDRequest struct {
	ID          string `param:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    int    `json:"priority" validate:"required,min=1,max=3"`
}

type TaskUpdateStatusWithIDRequest struct {
	ID     string `param:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type TaskDeleteRequest struct {
	ID string `param:"id" validate:"required"`
}

type EmptyRequest struct{}
