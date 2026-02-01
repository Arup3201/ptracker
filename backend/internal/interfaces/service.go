package interfaces

import (
	"context"

	"github.com/ptracker/internal/domain"
)

type AuthService interface {
	RedirectLogin(ctx context.Context) (string, error)
	Callback(ctx context.Context,
		state, code string,
		userAgent, device, ipAddress string) (*domain.Session, error)
	Authenticate(ctx context.Context,
		sessionId string) (string, error)
	Refresh(ctx context.Context,
		sessionId string) error
	Logout(ctx context.Context,
		sessionId string) error
}

type LimiterService interface {
	IsAllowed(ctx context.Context, userId string) (bool, error)
	GetTokens(ctx context.Context, userId string) (int, error)
	GetCapacity(ctx context.Context) int
	GetRetryTime(ctx context.Context, userId string) (int, error)
}

type ProjectService interface {
	CreateProject(ctx context.Context, name string,
		description, skills *string,
		owner string) (string, error)
	ListProjects(ctx context.Context, userId string) ([]*domain.PrivateProjectListed, error)
	GetPrivateProject(ctx context.Context,
		projectId, userId string) (*domain.ProjectDetail, error)
	GetProjectMembers(ctx context.Context,
		projectId, userId string) ([]*domain.Member, error)
	ListJoinRequests(ctx context.Context,
		projectId, userId string) ([]*domain.JoinRequestListed, error)
	RespondToJoinRequests(ctx context.Context,
		projectId, ownerId, requestorId, joinStatus string) error
}

type PublicService interface {
	ListPublicProjects(ctx context.Context,
		userId string) ([]*domain.PublicProjectListed, error)
	GetPublicProject(ctx context.Context,
		projectId string) (*domain.PublicProjectSummary, error)
	JoinProject(ctx context.Context,
		projectId, userId string) error
}

type TaskService interface {
	CreateTask(ctx context.Context,
		projectId, title string,
		description *string,
		userId string) (string, error)
	ListTasks(ctx context.Context,
		projectId, userId string) ([]*domain.TaskListed, error)
	GetTask(ctx context.Context,
		projectId, taskId, userId string) (*domain.Task, error)
	GetTaskAssignees(ctx context.Context,
		projectId, taskId, userId string) ([]*domain.Assignee, error)
}
