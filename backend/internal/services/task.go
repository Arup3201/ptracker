package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
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
	assignees []string,
	status, userId string) (string, []string, error) {

	permitted, err := s.permissionService.CanCreateTask(ctx, projectId, userId)
	if err != nil {
		return "", nil, fmt.Errorf("permission service can create task: %w", err)
	}

	if !permitted {
		return "", nil, apierr.ErrForbidden
	}

	if strings.Trim(title, " ") == "" {
		return "", nil, apierr.ErrInvalidValue
	}

	if !slices.Contains([]string{
		domain.TASK_STATUS_UNASSIGNED,
		domain.TASK_STATUS_ONGOING,
		domain.TASK_STATUS_COMPLETED,
		domain.TASK_STATUS_ABANDONED,
	}, status) {
		return "", nil, apierr.ErrInvalidValue
	}

	taskId, err := s.store.Task().Create(ctx,
		projectId, title, description, status)
	if err != nil {
		return "", nil, fmt.Errorf("store task create: %w", err)
	}

	warnings := s.AddAssignees(ctx, projectId, taskId, assignees)

	return taskId, warnings, nil
}

func (s *taskService) AddAssignees(ctx context.Context,
	projectId, taskId string,
	assignees []string) (warnings []string) {
	if len(assignees) == 0 {
		return
	}

	for _, assignee := range assignees {
		_, err := s.store.User().Get(ctx, assignee)
		if err != nil {
			switch err {
			case apierr.ErrNotFound:
				warnings = append(warnings, fmt.Sprintf("not a valid user %s", assignee))
			default:
				warnings = append(warnings, "failed to fetch user from database")
			}
			continue
		}

		err = s.store.Assignee().Create(ctx, projectId, taskId, assignee)
		if err != nil {
			warnings = append(warnings, "failed to create assignee in the database")
		}
	}

	return
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
