package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
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

	taskStatus := domain.TASK_STATUS_UNASSIGNED
	taskId, err := s.store.Task().Create(ctx,
		projectId, title, description, taskStatus)
	if err != nil {
		return "", fmt.Errorf("store task create: %w", err)
	}

	return taskId, nil
}

func (s *taskService) ListTasks(ctx context.Context,
	projectId, userId string) ([]*domain.TaskListed, error) {

	permitted, err := s.permissionService.CanSeeTasks(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can create task: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	tasks, err := s.store.List().Tasks(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("store list tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) GetTask(ctx context.Context,
	projectId, taskId, userId string) (*domain.Task, error) {

	permitted, err := s.permissionService.CanAccessTask(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can create task: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	task, err := s.store.Task().Get(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("store task get: %w", err)
	}

	return task, err
}

func (s *taskService) GetTaskAssignees(ctx context.Context,
	projectId, taskId, userId string) ([]*domain.Assignee, error) {

	permitted, err := s.permissionService.CanAccessTask(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can access task: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	assignees, err := s.store.List().Assignees(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("store list assignees: %w", err)
	}

	return assignees, nil
}
