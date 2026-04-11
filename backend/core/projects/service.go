package projects

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ptracker/core"
	"github.com/ptracker/core/members"
	"gorm.io/gorm"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Skills      *string   `json:"skills"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectSummary struct {
	Project
	UnassignedTasks int64 `json:"unassigned_tasks"`
	OngoingTasks    int64 `json:"ongoing_tasks"`
	CompletedTasks  int64 `json:"completed_tasks"`
	AbandonedTasks  int64 `json:"abandoned_tasks"`
}

type ProjectService struct {
	txManager   *core.TxManager
	projectRepo *ProjectRepository
	memberRepo  *members.MemberRepository
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) CreateProject(ctx context.Context,
	name string,
	description, skills *string,
	ownerID string) (string, error) {

	var err error
	var projectID string

	if strings.Trim(name, " ") == "" {
		return "", core.ErrInvalidValue
	}

	err = s.txManager.WithTx(func(tx *gorm.DB) error {

		projectRepo := s.projectRepo.WithTx(tx)
		memberRepo := s.memberRepo.WithTx(tx)

		projectID, err = projectRepo.Create(ctx, name, description, skills, ownerID)
		if err != nil {
			return fmt.Errorf("project repository create: %w", err)
		}

		err = memberRepo.Create(ctx, projectID, ownerID, core.ROLE_OWNER)
		if err != nil {
			return fmt.Errorf("member repository create: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("txManager WithTx: %w", err)
	}

	return projectID, nil
}

func (s *ProjectService) Get(ctx context.Context,
	projectID, userID string) (*ProjectSummary, error) {

	var err error

	err = core.NeedsToBeAMember(ctx, s.memberRepo, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("needs to be a member: %w", err)
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project repository get: %w", err)
	}

	myProject := ProjectSummary{
		Project: Project{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Skills:      project.Skills,
			OwnerID:     project.OwnerID,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		},
		UnassignedTasks: project.UnassignedTasks,
		OngoingTasks:    project.OngoingTasks,
		CompletedTasks:  project.CompletedTasks,
		AbandonedTasks:  project.AbandonedTasks,
	}
	return &myProject, nil
}

func (s *ProjectService) MyProjects(ctx context.Context,
	userID string) ([]ProjectSummary, error) {

	projects, err := s.projectRepo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("project repository list: %w", err)
	}

	myProjects := []ProjectSummary{}
	for _, p := range projects {
		myProjects = append(myProjects, ProjectSummary{
			Project: Project{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Skills:      p.Skills,
				OwnerID:     p.OwnerID,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			},
			UnassignedTasks: p.UnassignedTasks,
			OngoingTasks:    p.OngoingTasks,
			CompletedTasks:  p.CompletedTasks,
			AbandonedTasks:  p.AbandonedTasks,
		})
	}

	return myProjects, nil
}

func (s *ProjectService) RecentlyCreated(ctx context.Context,
	userID string) ([]ProjectSummary, error) {

	// pick last 10 recently created projects in descending order of their creation time
	rows, err := s.projectRepo.RecentlyCreated(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("list recently created projects: %w", err)
	}

	projects := []ProjectSummary{}
	for _, p := range rows {
		projects = append(projects, ProjectSummary{
			Project: Project{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Skills:      p.Skills,
				OwnerID:     p.OwnerID,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			},
			UnassignedTasks: p.UnassignedTasks,
			OngoingTasks:    p.OngoingTasks,
			CompletedTasks:  p.CompletedTasks,
			AbandonedTasks:  p.AbandonedTasks,
		})
	}

	return projects, nil
}

func (s *ProjectService) RecentlyJoined(ctx context.Context,
	userID string) ([]ProjectSummary, error) {

	// pick last 10 recently joined projects in descending order of their joining time
	rows, err := s.projectRepo.RecentlyJoined(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("list recently joined projects: %w", err)
	}

	projects := []ProjectSummary{}
	for _, p := range rows {
		projects = append(projects, ProjectSummary{
			Project: Project{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Skills:      p.Skills,
				OwnerID:     p.OwnerID,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			},
			UnassignedTasks: p.UnassignedTasks,
			OngoingTasks:    p.OngoingTasks,
			CompletedTasks:  p.CompletedTasks,
			AbandonedTasks:  p.AbandonedTasks,
		})
	}

	return projects, nil
}
