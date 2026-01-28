package stores

import (
	"context"

	"github.com/ptracker/domain"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*domain.User, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
	All(ctx context.Context, userId string) ([]*domain.ListedProject, error)
	Get(ctx context.Context, id string) (*domain.ListedProject, error)
}

type RoleRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
	Get(ctx context.Context, projectId, userId string) (string, error)
	CountMembers(ctx context.Context, projectId string) (int, error)
}

type Store interface {
	WithTx(ctx context.Context, fn func(txStore Store) error) error
	User() UserRepository
	Project() ProjectRepository
	Role() RoleRepository
}
