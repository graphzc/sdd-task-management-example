package router

import (
	"github.com/graphzc/sdd-task-management-example/internal/handlers"
	"github.com/graphzc/sdd-task-management-example/internal/middlewares"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo           *echo.Echo
	handlers       *handlers.Handlers
	authMiddleware middlewares.AuthMiddleware
}

func NewRouter(echo *echo.Echo, handlers *handlers.Handlers, authMiddleware middlewares.AuthMiddleware) *Router {
	return &Router{
		echo:           echo,
		handlers:       handlers,
		authMiddleware: authMiddleware,
	}
}
