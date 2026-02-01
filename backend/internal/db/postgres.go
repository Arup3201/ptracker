package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func ConnectPostgres(connString string) (*sql.DB, error) {
	var err error
	if connString == "" {
		return nil, fmt.Errorf("connect postgres: missing connection string")
	}
	var pgDb *sql.DB
	pgDb, err = sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	return pgDb, nil
}

// migrations: file source for migrations. It is path/to/file/or/folder when path is relative and /path/to/file/or/folder when absolute.
// https://github.com/golang-migrate/migrate/tree/master/source/file
func Migrate(migrations string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres migrate: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrations,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("postgres migrate: %s", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("postgres migrate: %s", err)
	}

	return nil
}
