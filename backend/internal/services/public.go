package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type publicService struct {
	store interfaces.Store
}

func NewPublicService(store interfaces.Store) interfaces.PublicService {
	return &publicService{
		store: store,
	}
}

func (s *publicService) ListPublicProjects(ctx context.Context, userId string) ([]*domain.PublicProjectListed, error) {
	projects, err := s.store.List().PublicProjects(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("store list public projects: %w", err)
	}

	return projects, nil
}

func (s *publicService) GetPublicProject(ctx context.Context,
	projectId, userId string) (*domain.PublicProjectSummary, error) {

	project, err := s.store.Public().Get(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store public get: %w", err)
	}

	owner, err := s.store.User().Get(ctx, project.Owner.UserId)
	if err != nil {
		return nil, fmt.Errorf("store user get: %w", err)
	}

	role, err := s.store.Role().Get(ctx, projectId, project.Owner.UserId)
	if err != nil {
		return nil, fmt.Errorf("store role get: %w", err)
	}

	project.Owner = &domain.Member{
		UserId:      project.Owner.UserId,
		Username:    owner.Username,
		DisplayName: owner.DisplayName,
		Email:       owner.Email,
		AvatarURL:   owner.AvatarURL,
		IsActive:    owner.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	joinStatus, err := s.store.JoinRequest().Get(ctx, projectId, userId)
	if err == apierr.ErrNotFound {
		project.JoinStatus = "Not Requested"
	} else if err != nil {
		return nil, fmt.Errorf("store join request get: %w", err)
	} else {
		project.JoinStatus = joinStatus
	}

	return project, nil
}

func (s *publicService) JoinProject(ctx context.Context,
	projectId, userId string) error {

	joinStatus := domain.JOIN_STATUS_PENDING

	err := s.store.JoinRequest().Create(ctx, projectId, userId, joinStatus)
	if err != nil {
		if errors.Is(err, apierr.ErrDuplicate) {
			return apierr.ErrDuplicate
		}
		return fmt.Errorf("store join request create: %w", err)
	}

	return nil
}
