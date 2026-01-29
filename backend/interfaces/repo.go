package interfaces

import (
	"context"

	"github.com/ptracker/domain"
)

type UserRepository interface {
	Create(ctx context.Context,
		idpSubject, idpProvider, username, email string,
		displayName, avatarUrl *string) (string, error)
	Get(ctx context.Context, id string) (*domain.User, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
	Get(ctx context.Context, id string) (*domain.ProjectSummary, error)
	Delete(ctx context.Context, id string) error
}

type RoleRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
	Get(ctx context.Context, projectId, userId string) (string, error)
	CountMembers(ctx context.Context, projectId string) (int, error)
}

type ListRepository interface {
	PrivateProjects(ctx context.Context, userId string) ([]*domain.PrivateProjectListed, error)
	Members(ctx context.Context, projectId string) ([]*domain.Member, error)
	PublicProjects(ctx context.Context, userId string) ([]*domain.PublicProjectListed, error)
}
