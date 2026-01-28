package fixtures

import (
	"context"

	"github.com/ptracker/stores"
)

type Fixtures struct {
	ctx   context.Context
	store stores.Store
}

func New(ctx context.Context, store stores.Store) *Fixtures {
	return &Fixtures{
		ctx:   ctx,
		store: store,
	}
}
