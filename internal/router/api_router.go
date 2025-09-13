package router

import (
	"net/http"

	"github.com/graphzc/sdd-task-management-example/internal/utils/echoutil"
)

func (r *Router) RegisterAPIRoutes() {
	// Health check
	r.echo.GET("/health", echoutil.WrapWithStatus(r.handlers.Common.HealthCheck, http.StatusOK))

	v1Public := r.echo.Group("/api/v1")

	// Auth routes
	authGroup := v1Public.Group("/auth")
	{
		authGroup.POST("/register", echoutil.WrapWithStatus(r.handlers.Auth.Register, http.StatusCreated))
		authGroup.POST("/login", echoutil.WrapWithStatus(r.handlers.Auth.Login, http.StatusOK))
	}
}
