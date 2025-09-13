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

	// Protected routes
	v1Protected := v1Public.Group("", r.authMiddleware.Middleware)

	// Task routes
	taskGroup := v1Protected.Group("/tasks")
	{
		taskGroup.POST("", echoutil.WrapWithStatus(r.handlers.Task.CreateTaskWrapped, http.StatusCreated))
		taskGroup.GET("", echoutil.WrapWithStatus(r.handlers.Task.GetTasksByUserIDWrapped, http.StatusOK))
		taskGroup.GET("/:id", echoutil.WrapWithStatus(r.handlers.Task.GetTaskByIDWrapped, http.StatusOK))
		taskGroup.PUT("/:id", echoutil.WrapWithStatus(r.handlers.Task.UpdateTaskByIDWrapped, http.StatusOK))
		taskGroup.PATCH("/:id/status", echoutil.WrapWithStatus(r.handlers.Task.UpdateTaskStatusByIDWrapped, http.StatusOK))
		taskGroup.DELETE("/:id", echoutil.WrapWithStatus(r.handlers.Task.DeleteTaskByIDWrapped, http.StatusOK))
	}
}
