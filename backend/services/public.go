package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
)

type publicService struct {
	store interfaces.Store
}

func NewPublicService(store interfaces.Store) *publicService {
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

func (s *publicService) JoinProject(ctx context.Context,
	projectId, userId string) error {

	joinStatus := "Pending"
	userRole := domain.ROLE_MEMBER

	err := s.store.WithTx(ctx, func(txStore interfaces.Store) error {
		var err error

		err = txStore.JoinRequest().Create(ctx, projectId, userId, joinStatus)
		if err != nil {
			if errors.Is(err, apierr.ErrDuplicate) {
				return apierr.ErrDuplicate
			}
			return fmt.Errorf("store join request create: %w", err)
		}

		err = txStore.Role().Create(ctx, projectId, userId, userRole)
		if err != nil {
			return fmt.Errorf("store role create: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
