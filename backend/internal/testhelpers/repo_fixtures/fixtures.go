package repo_fixtures

import (
	"context"
	"database/sql"
)

type Fixtures struct {
	ctx context.Context
	db  *sql.DB
}

func New(ctx context.Context, db *sql.DB) *Fixtures {
	return &Fixtures{
		ctx: ctx,
		db:  db,
	}
}
