package main

import (
	"github.com/graphzc/sdd-task-management-example/cmd/api/di"
	"github.com/rs/zerolog/log"
)

func main() {
	server := di.InitializeAPI()
	if err := server.Start(); err != nil {
		log.Panic().
			Err(err).
			Msg("Failed to start server")
	}
}
