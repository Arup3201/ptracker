package services

import "context"

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
}

type RoleRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
}

type Store interface {
	WithTx(ctx context.Context, fn func(txStore Store) error) error
	Project() ProjectRepository
	Role() RoleRepository
}
