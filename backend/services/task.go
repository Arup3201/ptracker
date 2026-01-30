package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ptracker/apierr"
	"github.com/ptracker/interfaces"
)

type taskService struct {
	store             interfaces.Store
	permissionService *ProjectPermissionService
}

func NewTaskService(store interfaces.Store) *taskService {
	permissionService := &ProjectPermissionService{
		store: store,
	}
	return &taskService{
		store:             store,
		permissionService: permissionService,
	}
}

func (s *taskService) CreateTask(ctx context.Context,
	projectId, title string,
	description *string,
	userId string) (string, error) {

	permitted, err := s.permissionService.CanCreateTask(ctx, projectId, userId)
	if err != nil {
		return "", fmt.Errorf("permission service can create task: %w", err)
	}

	if !permitted {
		return "", apierr.ErrForbidden
	}

	if strings.Trim(title, " ") == "" {
		return "", apierr.ErrInvalidValue
	}

	taskStatus := "Unassigned"
	taskId, err := s.store.Task().Create(ctx,
		projectId, title, description, taskStatus)
	if err != nil {
		return "", fmt.Errorf("store task create: %w", err)
	}

	return taskId, nil
}
