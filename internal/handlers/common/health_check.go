package common

import (
	"github.com/graphzc/sdd-task-management-example/internal/dto"
)

func (h *handler) HealthCheck(_ any) (dto.HealthCheckResponse, error) {
	return dto.HealthCheckResponse{
		Status: "ok",
	}, nil
}
