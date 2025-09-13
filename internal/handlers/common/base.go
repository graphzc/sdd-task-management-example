package common

import (
	"context"

	"github.com/graphzc/sdd-task-management-example/internal/dto"
)

type Handler interface {
	HealthCheck(ctx context.Context, _ any) (dto.HealthCheckResponse, error)
}

type handler struct{}

// @WireSet("Handler")
func New() Handler {
	return &handler{}
}
