package services

import (
	"context"
	"fmt"

	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
)

type ProjectPermissionService struct {
	store stores.Store
}

func (s *ProjectPermissionService) CanAccess(ctx context.Context,
	projectId, userId string) (bool, error) {
	role, err := s.store.Role().Get(ctx, projectId, userId)
	if err != nil {
		return false, fmt.Errorf("role store get: %w", err)
	}

	if role == domain.ROLE_OWNER || role == domain.ROLE_MEMBER {
		return true, nil
	}

	return false, nil
}
