package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type ProjectPermissionService struct {
	store interfaces.Store
}

func (s *ProjectPermissionService) CanAccess(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER || role.Role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanSeeMembers(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER || role.Role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanRespondToJoinRequests(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanCreateTask(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanSeeTasks(ctx context.Context,
	projectId, userId string) (bool, error) {

	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER || role.Role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanAccessTask(ctx context.Context,
	projectId, userId string) (bool, error) {

	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER || role.Role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanUpdateTask(ctx context.Context,
	projectId, taskId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER {
		return true, nil
	}

	exist, err := s.store.Assignee().Get(ctx, projectId, taskId, userId)
	if err == apierr.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("store assignee get: %w", err)
	}

	if exist {
		return true, nil
	}

	return false, nil
}

func (s *ProjectPermissionService) CanCommentOnTask(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role.Role == domain.ROLE_OWNER || role.Role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}
