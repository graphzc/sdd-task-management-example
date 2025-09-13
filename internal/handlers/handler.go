package handlers

import (
	"github.com/graphzc/sdd-task-management-example/internal/handlers/auth"
	"github.com/graphzc/sdd-task-management-example/internal/handlers/common"
	"github.com/graphzc/sdd-task-management-example/internal/handlers/task"
)

type Handlers struct {
	Common common.Handler
	Auth   auth.Handler
	Task   task.Handler
}

// @WireSet("Handler")
func NewHandlers(
	commonHandler common.Handler,
	authHandler auth.Handler,
	taskHandler task.Handler,
) *Handlers {
	return &Handlers{
		Common: commonHandler,
		Auth:   authHandler,
		Task:   taskHandler,
	}
}
