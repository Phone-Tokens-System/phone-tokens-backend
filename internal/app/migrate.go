package app

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

// MigrateDB applies SQL migrations located in migrationsDir using goose.
func MigrateDB(db *sql.DB, migrationsDir string) error {
	goose.SetDialect("postgres")
	return goose.Up(db, migrationsDir)
}
