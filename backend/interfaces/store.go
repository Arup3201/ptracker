package interfaces

import (
	"context"
)

type Store interface {
	WithTx(ctx context.Context, fn func(txStore Store) error) error
	User() UserRepository
	Project() ProjectRepository
	Role() RoleRepository
	List() ListRepository
}
