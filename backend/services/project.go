package services

import (
	"context"
	"fmt"

	"github.com/ptracker/domain"
)

type projectService struct {
	store Store
}

func NewProjectService(store Store) *projectService {
	return &projectService{
		store: store,
	}
}

func (s *projectService) CreateProject(ctx context.Context, name string,
	description, skills *string,
	owner string) (string, error) {
	var projectId string

	err := s.store.WithTx(ctx, func(txStore Store) error {
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

func (s *projectService) ListProjects(ctx context.Context, userId string) ([]domain.ProjectSummary, error) {
	projects, err := s.store.Project().All(ctx, userId)
	if err != nil {
		return projects, fmt.Errorf("store project all: %w", err)
	}

	return projects, nil
}
