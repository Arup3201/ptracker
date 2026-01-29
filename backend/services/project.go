package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
)

type projectService struct {
	store             stores.Store
	projectPermission *ProjectPermissionService
}

func NewProjectService(store stores.Store) *projectService {
	permissionService := &ProjectPermissionService{
		store: store,
	}
	return &projectService{
		store:             store,
		projectPermission: permissionService,
	}
}

func (s *projectService) CreateProject(ctx context.Context, name string,
	description, skills *string,
	owner string) (string, error) {
	var projectId string

	if strings.Trim(name, " ") == "" {
		return "", apierr.ErrInvalidValue
	}

	err := s.store.WithTx(ctx, func(txStore stores.Store) error {
		var err error

		projectId, err = s.store.Project().Create(ctx, name, description, skills, owner)
		if err != nil {
			return fmt.Errorf("store project create: %w", err)
		}

		err = s.store.Role().Create(ctx, projectId, owner, domain.ROLE_OWNER)
		if err != nil {
			return fmt.Errorf("store role create: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("store WithTx: %w", err)
	}

	return projectId, nil
}

func (s *projectService) ListProjects(ctx context.Context, userId string) ([]*domain.PrivateProjectListed, error) {
	projects, err := s.store.List().PrivateProjects(ctx, userId)
	if err != nil {
		return projects, fmt.Errorf("store project all: %w", err)
	}

	return projects, nil
}

func (s *projectService) GetPrivateProject(ctx context.Context,
	projectId, userId string) (*domain.ProjectDetail, error) {

	permitted, err := s.projectPermission.CanAccess(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can access: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	project, err := s.store.Project().Get(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store project all: %w", err)
	}

	userRole, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("store role get: %w", err)
	}

	owner, err := s.store.User().Get(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("store user get: %w", err)
	}

	memberCount, err := s.store.Role().CountMembers(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store role count members: %w", err)
	}

	return &domain.ProjectDetail{
		ProjectSummary: domain.ProjectSummary{
			PrivateProject: domain.PrivateProject{
				Id:          project.Id,
				Name:        project.Name,
				Description: project.Description,
				Skills:      project.Skills,
				CreatedAt:   project.CreatedAt,
				UpdatedAt:   project.UpdatedAt,
			},
			UnassignedTasks: project.UnassignedTasks,
			OngoingTasks:    project.OngoingTasks,
			CompletedTasks:  project.CompletedTasks,
			AbandonedTasks:  project.AbandonedTasks,
		},
		Role:        userRole,
		MemberCount: memberCount,
		Owner: &domain.Member{
			Id:          owner.Id,
			Username:    owner.Username,
			DisplayName: owner.DisplayName,
			Email:       owner.Email,
			AvatarURL:   owner.AvatarURL,
			IsActive:    owner.IsActive,
			CreatedAt:   owner.CreatedAt,
			UpdatedAt:   owner.UpdatedAt,
		},
	}, nil
}
