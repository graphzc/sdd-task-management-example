package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/graphzc/sdd-task-management-example/internal/domain/enums"
)

type Model struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	UserID      uuid.UUID          `json:"userId" db:"user_id"`
	Title       string             `json:"title" db:"title"`
	Description string             `json:"description" db:"description"`
	Priority    enums.TaskPriority `json:"priority" db:"priority"`
	Status      enums.TaskStatus   `json:"status" db:"status"`
	CreatedAt   time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" db:"updated_at"`
}
