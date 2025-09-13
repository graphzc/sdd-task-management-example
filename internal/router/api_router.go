package router

import (
	"net/http"

	"github.com/graphzc/sdd-task-management-example/internal/utils/echoutil"
)

func (r *Router) RegisterAPIRoutes() {
	// Health check
	r.echo.GET("/health", echoutil.WrapWithStatus(r.handlers.Common.HealthCheck, http.StatusOK))
}
