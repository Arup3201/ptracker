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

func (s *publicService) ListPublicProjects(ctx context.Context,
	userId string) ([]domain.ProjectPreview, error) {
	projects, err := s.store.List().PublicProjects(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("store list public projects: %w", err)
	}

	return projects, nil
}

func (s *publicService) GetPublicProject(ctx context.Context,
	projectId, userId string) (*domain.ProjectPublicDetail, error) {

	project, err := s.store.Public().Get(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store public get: %w", err)
	}

	return &project, nil
}

func (s *publicService) GetJoinStatus(ctx context.Context,
	projectId, userId string) (string, error) {
	joinStatus, err := s.store.JoinRequest().Status(ctx, projectId, userId)
	if err == apierr.ErrNotFound {
		joinStatus = "Not Requested"
	} else if err != nil {
		return "", fmt.Errorf("store join request get: %w", err)
	}

	return joinStatus, nil
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
