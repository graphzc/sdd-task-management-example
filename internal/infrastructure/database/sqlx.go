package database

import (
	"context"

	"github.com/graphzc/sdd-task-management-example/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// @WireSet("Infrastructure")
func NewSQLXClient(ctx context.Context, config *config.Config) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, config.Database.Driver, config.Database.URI)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Failed to connect to the database")
	}

	log.Info().
		Msg("Connected to the database")

	return db
}
