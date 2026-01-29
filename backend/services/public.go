package services

import (
	"context"
	"fmt"

	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
)

type publicService struct {
	store interfaces.Store
}

func (s *publicService) ListPublicProjects(ctx context.Context, userId string) ([]*domain.PublicProjectListed, error) {
	projects, err := s.store.List().PublicProjects(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("store list public projects: %w", err)
	}

	return projects, nil
}
