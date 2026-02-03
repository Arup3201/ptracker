package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDatabase(driver, dataSource string) (*sql.DB, error) {
	var err error

	if dataSource == "" {
		return nil, fmt.Errorf("missing connection string")
	}

	var db *sql.DB
	db, err = sql.Open(driver, dataSource)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}
