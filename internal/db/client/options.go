package db

import (
	"table-app/internal/db/migration"
	"table-app/internal/log"
)

type Option func(db *Client)

func WithMigrationRunner(migrationDir string, logger log.Logger) Option {
	return func(db *Client) {
		db.migrationRunner = migration.NewRunner(logger, migration.DialectPostgreSQL, migrationDir)
	}
}
