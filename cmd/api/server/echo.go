package server

import (
	"fmt"

	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/graphzc/sdd-task-management-example/internal/handlers"
	"github.com/graphzc/sdd-task-management-example/internal/middlewares"
	"github.com/graphzc/sdd-task-management-example/internal/router"
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
	"github.com/graphzc/sdd-task-management-example/internal/utils/validator"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	config         *config.Config
	handlers       *handlers.Handlers
	authMiddleware middlewares.AuthMiddleware
}

func NewEchoServer(
	config *config.Config,
	handlers *handlers.Handlers,
	authMiddleware middlewares.AuthMiddleware,
) *EchoServer {
	return &EchoServer{
		config:         config,
		handlers:       handlers,
		authMiddleware: authMiddleware,
	}
}

func (s *EchoServer) Start() error {
	e := echo.New()

	e.Validator = validator.NewValidator()

	e.HTTPErrorHandler = servererr.EchoHTTPErrorHandler

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     s.config.CORS.AllowOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	router := router.NewRouter(e, s.handlers, s.authMiddleware)

	router.RegisterAPIRoutes()

	return e.Start(fmt.Sprintf(":%s", s.config.Port))
}
