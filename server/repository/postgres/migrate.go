package postgres

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func applyMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "repository/postgres/psql_migrations",
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
