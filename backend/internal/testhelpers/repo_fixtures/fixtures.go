package repo_fixtures

import (
	"context"

	"github.com/ptracker/internal/interfaces"
)

type Fixtures struct {
	ctx context.Context
	db  interfaces.Execer
}

func New(ctx context.Context, db interfaces.Execer) *Fixtures {
	return &Fixtures{
		ctx: ctx,
		db:  db,
	}
}
