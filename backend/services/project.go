package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
)

type projectService struct {
	store             interfaces.Store
	projectPermission *ProjectPermissionService
}

func NewProjectService(store interfaces.Store) *projectService {
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

	err := s.store.WithTx(ctx, func(txStore interfaces.Store) error {
		var err error

		projectId, err = txStore.Project().Create(ctx, name, description, skills, owner)
		if err != nil {
			return fmt.Errorf("store project create: %w", err)
		}

		err = txStore.Role().Create(ctx, projectId, owner, domain.ROLE_OWNER)
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
		ProjectSummary: &domain.ProjectSummary{
			Project: &domain.Project{
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
		Role:        userRole.Role,
		MemberCount: memberCount,
		Owner: &domain.Member{
			Id:          owner.Id,
			Username:    owner.Username,
			DisplayName: owner.DisplayName,
			Email:       owner.Email,
			AvatarURL:   owner.AvatarURL,
			IsActive:    owner.IsActive,
			CreatedAt:   userRole.CreatedAt,
			UpdatedAt:   userRole.UpdatedAt,
		},
	}, nil
}

func (s *projectService) GetProjectMembers(ctx context.Context,
	projectId, userId string) ([]*domain.Member, error) {

	permitted, err := s.projectPermission.CanSeeMembers(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can access: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	members, err := s.store.List().Members(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store list members: %w", err)
	}

	return members, nil
}

func (s *projectService) ListJoinRequests(ctx context.Context,
	projectId, userId string) ([]*domain.JoinRequestListed, error) {

	permitted, err := s.projectPermission.CanSeeMembers(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can access: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	joinRequests, err := s.store.List().JoinRequests(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store list join requests: %w", err)
	}

	return joinRequests, nil
}

func (s *projectService) AccessOrRejectJoinRequest(ctx context.Context,
	projectId, ownerId, requestorId, joinStatus string) error {

	permitted, err := s.projectPermission.CanRespondToJoinRequests(ctx, projectId, ownerId)
	if err != nil {
		return fmt.Errorf("permission service can access: %w", err)
	}
	if !permitted {
		return apierr.ErrForbidden
	}

	if !slices.Contains([]string{
		domain.JOIN_STATUS_PENDING,
		domain.JOIN_STATUS_ACCEPTED,
		domain.JOIN_STATUS_REJECTED,
	}, joinStatus) {
		return apierr.ErrInvalidValue
	}

	userRole := domain.ROLE_MEMBER

	status, err := s.store.JoinRequest().Get(ctx, projectId, requestorId)
	if err != nil {
		return fmt.Errorf("store join request get: %w", err)
	}

	if status == joinStatus {
		return apierr.ErrInvalidValue
	}

	// Y Rejected -> Pending
	// X Rejected -> Accepted
	// X Accepted -> Pending|Rejected

	if (status == domain.JOIN_STATUS_REJECTED &&
		slices.Contains([]string{
			domain.JOIN_STATUS_ACCEPTED,
		}, joinStatus)) ||
		(status == domain.JOIN_STATUS_ACCEPTED &&
			slices.Contains([]string{
				domain.JOIN_STATUS_PENDING,
				domain.JOIN_STATUS_REJECTED,
			}, joinStatus)) {
		return apierr.ErrInvalidValue
	}

	err = s.store.WithTx(ctx, func(txStore interfaces.Store) error {
		var err error

		err = txStore.JoinRequest().Update(ctx, projectId, requestorId, joinStatus)
		if err != nil {
			return fmt.Errorf("store join request update: %w", err)
		}

		if joinStatus == domain.JOIN_STATUS_ACCEPTED {
			err = txStore.Role().Create(ctx, projectId, requestorId, userRole)
			if err != nil {
				return fmt.Errorf("store role create: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
