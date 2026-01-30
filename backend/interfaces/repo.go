package interfaces

import (
	"context"
	"time"

	"github.com/ptracker/domain"
)

type SessionRepository interface {
	Create(ctx context.Context,
		userId string,
		encryptedToken []byte,
		userAgent, ipAddress, deviceName string,
		expireAt time.Time) (string, error)
	Get(ctx context.Context, id string) (*domain.Session, error)
	Revoke(ctx context.Context, id string) error
	Update(ctx context.Context, id string,
		refreshTokenEncrypted []byte,
		expiresAt time.Time) error
}

type UserRepository interface {
	Create(ctx context.Context,
		idpSubject, idpProvider, username, email string,
		displayName, avatarUrl *string) (string, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	GetBySubject(ctx context.Context,
		idpSubject, idpProvider string) (*domain.User, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
	Get(ctx context.Context, id string) (*domain.ProjectSummary, error)
	Delete(ctx context.Context, id string) error
}

type TaskRepository interface {
	Create(ctx context.Context, projectId, title string,
		description *string,
		status string) (string, error)
	Get(ctx context.Context, id string) (*domain.Task, error)
}

type RoleRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
	Get(ctx context.Context, projectId, userId string) (*domain.Role, error)
	CountMembers(ctx context.Context, projectId string) (int, error)
}

type ListRepository interface {
	PrivateProjects(ctx context.Context, userId string) ([]*domain.PrivateProjectListed, error)
	Tasks(ctx context.Context, projectId string) ([]*domain.TaskListed, error)
	Assignees(ctx context.Context, taskId string) ([]*domain.Assignee, error)
	Members(ctx context.Context, projectId string) ([]*domain.Member, error)
	PublicProjects(ctx context.Context, userId string) ([]*domain.PublicProjectListed, error)
	JoinRequests(ctx context.Context, projectId string) ([]*domain.JoinRequestListed, error)
}

type JoinRequestRepository interface {
	Create(ctx context.Context, projectId, userId, joinStatus string) error
	Get(ctx context.Context, projectId, userId string) (string, error)
	Update(ctx context.Context, projectId, userId, joinStatus string) error
}

type PublicRepository interface {
	Get(ctx context.Context, projectId string) (*domain.PublicProjectSummary, error)
}
