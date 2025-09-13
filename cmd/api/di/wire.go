//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/graphzc/sdd-task-management-example/cmd/api/server"
)

func InitializeAPI() *server.EchoServer {
	wire.Build(
		ConfigSet,
		InfrastructureSet,
		HandlerSet,
		RepositorySet,
		ServiceSet,
		MiddlewareSet,
		server.NewEchoServer,
	)

	return &server.EchoServer{}
}
