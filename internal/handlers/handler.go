package handlers

import "github.com/graphzc/sdd-task-management-example/internal/handlers/common"

type Handlers struct {
	Common common.Handler
}

// @WireSet("Handler")
func NewHandlers(
	commonHandler common.Handler,
) *Handlers {
	return &Handlers{
		Common: commonHandler,
	}
}
