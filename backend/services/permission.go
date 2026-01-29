package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
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

	if role == domain.ROLE_OWNER || role == domain.ROLE_MEMBER {
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

	if role == domain.ROLE_OWNER || role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}
