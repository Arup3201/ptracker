package interfaces

import (
	"context"
	"time"

	"github.com/ptracker/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context,
		userId string,
		encryptedToken []byte,
		userAgent, ipAddress, deviceName string,
		expireAt time.Time) (string, error)
	Get(ctx context.Context, id string) (domain.Session, error)
	Revoke(ctx context.Context, id string) error
	Update(ctx context.Context, id string,
		refreshTokenEncrypted []byte,
		expiresAt time.Time) error
}

type UserRepository interface {
	Create(ctx context.Context,
		idpSubject, idpProvider, username, email string,
		displayName, avatarUrl *string) (string, error)
	Get(ctx context.Context, id string) (domain.User, error)
	GetBySubject(ctx context.Context,
		idpSubject, idpProvider string) (domain.User, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, title string,
		description, skills *string,
		owner string) (string, error)
	Get(ctx context.Context, id string) (domain.ProjectSummary, error)
	Delete(ctx context.Context, id string) error
}

type TaskRepository interface {
	Create(ctx context.Context, projectId,
		title, description, status string) (string, error)
	Get(ctx context.Context, id string) (domain.ProjectTaskItem, error)
	Update(ctx context.Context, id string,
		title, description, status *string) error
}

type MembershipRepository interface {
	Create(ctx context.Context, projectId, userId, role string) error
	Role(ctx context.Context, projectId, userId string) (string, error)
	CountMembers(ctx context.Context, projectId string) (int64, error)
}

type AssigneeRepository interface {
	Create(ctx context.Context, projectId, taskId, userId string) error
	Is(ctx context.Context, projectId, taskId, userId string) (bool, error)
	Delete(ctx context.Context, projectId, taskId, userId string) error
}

type ListRepository interface {
	PrivateProjects(ctx context.Context, userId string) ([]domain.ProjectSummary, error)
	Tasks(ctx context.Context, projectId string) ([]domain.ProjectTaskItem, error)
	Members(ctx context.Context, projectId string) ([]domain.Membership, error)
	PublicProjects(ctx context.Context, userId string) ([]domain.ProjectPreview, error)
	JoinRequests(ctx context.Context, projectId string) ([]domain.JoinRequest, error)
	Comments(ctx context.Context, projectId, taskId string) ([]domain.Comment, error)

	RecentlyCreatedProjects(ctx context.Context,
		userId string,
		n int) ([]domain.ProjectSummary, error)
	RecentlyJoinedProjects(ctx context.Context,
		userId string,
		n int) ([]domain.ProjectSummary, error)

	RecentlyAssignedTasks(ctx context.Context,
		userId string,
		n int) ([]domain.DashboardTaskItem, error)
	RecentlyUnassignedTasks(ctx context.Context,
		userId string,
		n int) ([]domain.DashboardTaskItem, error)
}

type JoinRequestRepository interface {
	Create(ctx context.Context, projectId, userId, joinStatus string) error
	Status(ctx context.Context, projectId, userId string) (string, error)
	Update(ctx context.Context, projectId, userId, joinStatus string) error
}

type PublicRepository interface {
	Get(ctx context.Context, projectId string) (domain.ProjectPublicDetail, error)
}

type CommentRepository interface {
	Create(ctx context.Context,
		projectId, taskId, userId string,
		comment string) (string, error)
}
