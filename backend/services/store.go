package services

import (
	"context"

	"github.com/ptracker/domain"
)

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
	All(ctx context.Context, userId string) ([]domain.ListedProject, error)
}

type RoleRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
}

type Store interface {
	WithTx(ctx context.Context, fn func(txStore Store) error) error
	Project() ProjectRepository
	Role() RoleRepository
}
