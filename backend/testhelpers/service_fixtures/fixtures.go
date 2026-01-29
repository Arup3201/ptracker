package service_fixtures

import (
	"context"

	"github.com/ptracker/interfaces"
)

type Fixtures struct {
	ctx   context.Context
	store interfaces.Store
}

func New(ctx context.Context, store interfaces.Store) *Fixtures {
	return &Fixtures{
		ctx:   ctx,
		store: store,
	}
}
