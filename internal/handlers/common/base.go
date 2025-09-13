package common

import (
	"github.com/graphzc/sdd-task-management-example/internal/dto"
)

type Handler interface {
	HealthCheck(_ any) (dto.HealthCheckResponse, error)
}

type handler struct {
}

// @WireSet("Handler")
func New() Handler {
	return &handler{}
}
