package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ptracker/internal/interfaces"
)

type database struct {
	db *sql.DB
}

func NewDatabase(driver, dataSource string) (interfaces.Execer, error) {
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

	return &database{
		db: db,
	}, nil
}

func (d *database) QueryContext(ctx context.Context, stmt string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, stmt, args)
}

func (d *database) QueryRowContext(ctx context.Context, stmt string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, stmt, args)
}

func (d *database) ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, stmt, args)
}
