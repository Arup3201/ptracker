package services

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/constants"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type taskService struct {
	store             interfaces.Store
	permissionService *ProjectPermissionService
	notifier          interfaces.Notifier
}

func NewTaskService(store interfaces.Store,
	notifier interfaces.Notifier) interfaces.TaskService {
	permissionService := &ProjectPermissionService{
		store: store,
	}
	return &taskService{
		store:             store,
		permissionService: permissionService,
		notifier:          notifier,
	}
}

func (s *taskService) CreateTask(ctx context.Context,
	projectId, userId string,
	title, description, status string,
	assignees []string) (string, []string, error) {

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

	project, _ := s.store.Project().Get(ctx, projectId)
	task, _ := s.store.Task().Get(ctx, taskId)

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

		role, err := s.store.Role().Get(ctx, projectId, assignee)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("failed to get user %s role", assignee))
			continue
		}

		if role.Role != domain.ROLE_MEMBER && role.Role != domain.ROLE_OWNER {
			warnings = append(warnings, fmt.Sprintf("user %s is not a member", assignee))
			continue
		}

		err = s.store.Assignee().Create(ctx, projectId, taskId, assignee)
		if err != nil {
			warnings = append(warnings, "failed to create assignee in the database")
		} else {
			var message domain.Message
			if project != nil && task != nil {
				message = domain.Message{
					Type: constants.ASSIGNEE_ADDED,
					Data: map[string]string{
						"project_id":   projectId,
						"task_id":      taskId,
						"project_name": project.Name,
						"task_name":    task.Title,
					},
				}

				err = s.notifier.Notify(ctx, assignee, message)
				if err != nil {
					fmt.Printf("[WARNING] notifier error: %s\n", err)
				}
			}
		}
	}

	return
}

func (s *taskService) RemoveAssignees(ctx context.Context,
	projectId, taskId string,
	assignees []string) (warnings []string) {
	if len(assignees) == 0 {
		return
	}

	project, _ := s.store.Project().Get(ctx, projectId)
	task, _ := s.store.Task().Get(ctx, taskId)

	for _, assignee := range assignees {
		exist, err := s.store.Assignee().Get(ctx, projectId, taskId, assignee)
		if err != nil {
			switch err {
			case apierr.ErrNotFound:
				warnings = append(warnings, fmt.Sprintf("not a valid user %s", assignee))
			default:
				warnings = append(warnings, "failed to fetch user from database")
			}
			continue
		}

		if !exist {
			warnings = append(warnings, fmt.Sprintf("user %s is not a valid assignee", assignee))
			continue
		}

		err = s.store.Assignee().Delete(ctx, projectId, taskId, assignee)
		if err != nil {
			warnings = append(warnings, "failed to create assignee in the database")
		} else {
			var message domain.Message
			if project != nil && task != nil {
				message = domain.Message{
					Type: constants.ASSIGNEE_REMOVED,
					Data: map[string]string{
						"project_id":   projectId,
						"task_id":      taskId,
						"project_name": project.Name,
						"task_name":    task.Title,
					},
				}
				err = s.notifier.Notify(ctx, assignee, message)
				if err != nil {
					fmt.Printf("[WARNING] notifier error: %s\n", err)
				}
			}

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

func (s *taskService) UpdateTask(ctx context.Context,
	projectId, taskId, userId string,
	title, description, status string,
	addedAssignees, removedAssignees []string) error {

	task, err := s.store.Task().Get(ctx, taskId)
	if err != nil {
		return fmt.Errorf("store task get: %w", err)
	}

	if title == task.Title &&
		(task.Description == nil && description == "") &&
		task.Status == status {
		return fmt.Errorf("no change")
	}

	if !slices.Contains([]string{
		domain.TASK_STATUS_UNASSIGNED,
		domain.TASK_STATUS_ONGOING,
		domain.TASK_STATUS_COMPLETED,
		domain.TASK_STATUS_ABANDONED,
	}, status) {
		return apierr.ErrInvalidValue
	}

	var hasUpdated = false
	taskTitle := task.Title

	if title != task.Title {
		permitted, err := s.permissionService.CanUpdateTask(ctx,
			projectId, taskId, userId)
		if err != nil {
			return fmt.Errorf("permission service can update task title: %w", err)
		}
		if !permitted {
			return apierr.ErrForbidden
		}

		task.Title = title
		hasUpdated = true
	}

	if task.Description != nil {
		desc := *task.Description
		if desc != description {
			permitted, err := s.permissionService.CanUpdateTask(ctx,
				projectId, taskId, userId)
			if err != nil {
				return fmt.Errorf("permission service can update task title: %w", err)
			}
			if !permitted {
				return apierr.ErrForbidden
			}

			*task.Description = description
			hasUpdated = true
		}
	} else if description != "" {
		permitted, err := s.permissionService.CanUpdateTask(ctx,
			projectId, taskId, userId)
		if err != nil {
			return fmt.Errorf("permission service can update task title: %w", err)
		}
		if !permitted {
			return apierr.ErrForbidden
		}

		task.Description = &description
		hasUpdated = true
	}

	if status != task.Status {
		permitted, err := s.permissionService.CanUpdateTask(ctx,
			projectId, taskId, userId)
		if err != nil {
			return fmt.Errorf("permission service can update task title: %w", err)
		}
		if !permitted {
			return apierr.ErrForbidden
		}

		task.Status = status
		hasUpdated = true
	}

	err = s.store.Task().Update(ctx, task)
	if err != nil {
		return fmt.Errorf("store task update: %w", err)
	}

	// Notifications

	if hasUpdated {
		message := domain.Message{
			Type: constants.TASK_UPDATED,
			Data: map[string]string{
				"task_id": taskId,
				"title":   taskTitle,
			},
		}

		project, _ := s.store.Project().Get(ctx, projectId)
		assignees, _ := s.store.List().Assignees(ctx, taskId)

		users := []string{}
		if userId != project.Owner {
			users = append(users, project.Owner)
		}
		for _, assignee := range assignees {
			if userId != assignee.UserId {
				users = append(users, assignee.UserId)
			}
		}
		s.notifier.BatchNotify(ctx, users, message)
	}

	s.AddAssignees(ctx, projectId, taskId, addedAssignees)
	s.RemoveAssignees(ctx, projectId, taskId, removedAssignees)

	return nil
}

func (s *taskService) AddComment(ctx context.Context,
	projectId, taskId, userId string,
	comment string) (string, error) {

	permitted, err := s.permissionService.CanCommentOnTask(ctx, projectId, userId)
	if err != nil {
		return "", fmt.Errorf("permission service can comment on task: %w", err)
	}

	if !permitted {
		return "", apierr.ErrForbidden
	}

	if strings.Trim(comment, " ") == "" {
		return "", apierr.ErrInvalidValue
	}

	id, err := s.store.Comment().Create(ctx, projectId, taskId, userId, comment)
	if err != nil {
		return "", fmt.Errorf("store comment create: %w", err)
	}

	// Notifications

	task, _ := s.store.Task().Get(ctx, taskId)
	message := domain.Message{
		Type: constants.COMMENT_ADDED,
		Data: map[string]string{
			"task_id": taskId,
			"title":   task.Title,
		},
	}

	project, _ := s.store.Project().Get(ctx, projectId)
	assignees, _ := s.store.List().Assignees(ctx, taskId)

	users := []string{}
	if userId != project.Owner {
		users = append(users, project.Owner)
	}
	for _, assignee := range assignees {
		if userId != assignee.UserId {
			users = append(users, assignee.UserId)
		}
	}
	s.notifier.BatchNotify(ctx, users, message)

	return id, nil
}

func (s *taskService) ListComments(ctx context.Context,
	projectId, taskId, userId string) ([]*domain.Comment, error) {

	permitted, err := s.permissionService.CanAccessTask(ctx, projectId, userId)
	if err != nil {
		return nil, fmt.Errorf("permission service can access task: %w", err)
	}

	if !permitted {
		return nil, apierr.ErrForbidden
	}

	comments, err := s.store.List().Comments(ctx, projectId, taskId)
	if err != nil {
		return nil, fmt.Errorf("store list comments: %w", err)
	}

	return comments, nil
}

func (s *taskService) AssignedTasks(ctx context.Context,
	userId string) ([]*domain.RecentTaskListed, error) {

	// pick last 10 recently joined projects in descending order of their joining time
	tasks, err := s.store.List().RecentlyAssignedTasks(ctx, userId, 10)
	if err != nil {
		return nil, fmt.Errorf("list recently joined tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) UnassignedTasks(ctx context.Context,
	userId string) ([]*domain.RecentTaskListed, error) {

	// pick last 10 recently joined projects in descending order of their joining time
	projects, err := s.store.List().RecentlyUnassignedTasks(ctx, userId, 10)
	if err != nil {
		return nil, fmt.Errorf("list recently joined projects: %w", err)
	}

	return projects, nil
}
