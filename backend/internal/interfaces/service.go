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
	Me(ctx context.Context,
		userId string) (*domain.User, error)
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
	ListProjects(ctx context.Context, userId string) ([]domain.ProjectSummary, error)
	GetPrivateProject(ctx context.Context,
		projectId, userId string) (*domain.ProjectDetail, error)
	GetProjectMembers(ctx context.Context,
		projectId, userId string) ([]domain.Membership, error)
	ListJoinRequests(ctx context.Context,
		projectId, userId string) ([]domain.JoinRequest, error)
	RespondToJoinRequests(ctx context.Context,
		projectId, ownerId, requestorId, joinStatus string) error

	ListRecentlyCreatedProjects(ctx context.Context,
		userId string) ([]domain.ProjectSummary, error)
	ListRecentlyJoinedProjects(ctx context.Context,
		userId string) ([]domain.ProjectSummary, error)
}

type PublicService interface {
	ListPublicProjects(ctx context.Context,
		userId string) ([]domain.ProjectPreview, error)
	GetPublicProject(ctx context.Context,
		projectId, userId string) (*domain.ProjectPublicDetail, error)
	GetJoinStatus(ctx context.Context,
		projectId, userId string) (string, error)
	JoinProject(ctx context.Context,
		projectId, userId string) error
}

type TaskService interface {
	CreateTask(ctx context.Context,
		projectId, userId string,
		title, description, status string,
		assignees []string) (string, []string, error)
	ListTasks(ctx context.Context,
		projectId, userId string) ([]domain.ProjectTaskItem, error)
	GetTask(ctx context.Context,
		projectId, taskId, userId string) (*domain.ProjectTaskItem, error)
	UpdateTask(ctx context.Context,
		projectId, taskId, userId string,
		title, description, status *string,
		addedAssignees, removedAssignees []string) error
	AddComment(ctx context.Context,
		projectId, taskId, userId string,
		comment string) (string, error)
	ListComments(ctx context.Context,
		projectId, taskId, userId string) ([]domain.Comment, error)

	AssignedTasks(ctx context.Context,
		userId string) ([]domain.DashboardTaskItem, error)
	UnassignedTasks(ctx context.Context,
		userId string) ([]domain.DashboardTaskItem, error)
}
