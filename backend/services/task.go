package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/ptracker/apierr"
	"github.com/ptracker/models"
)

type TaskServiceStore interface {
	WithTx(ctx context.Context, fn func(txStore TaskServiceStore) error) error
	Task() TaskStore
	Assignee() AssigneeStore
}

type taskService struct {
	TaskServiceStore
	access    *AccessControl
	ProjectId string
	UserId    string
}

func CreateTaskService(serviceStore TaskServiceStore,
	roleStore RoleStore,
	projectId, userId string) *taskService {

	return &taskService{
		TaskServiceStore: serviceStore,
		access: &AccessControl{
			roleStore: roleStore,
		},
		ProjectId: projectId,
		UserId:    userId,
	}
}

type TaskUpdates struct {
	Title       *string
	Description *string
	Status      *string
	Assignee    *string
}

func verifyTask(t *models.ProjectTask) error {
	taskStatus := []string{"Unassigned", "Ongoing", "Completed", "Abandoned"}
	if !slices.Contains(taskStatus, t.Status) {
		return apierr.ErrInvalidValue
	}

	if t.Title == "" {
		return apierr.ErrInvalidValue
	}

	return nil
}

func (ts *taskService) PatchTask(id string, taskChanges TaskUpdates) error {
	isOwner, err := ts.Access().IsOwner()
	if err != nil {
		return fmt.Errorf("access is owner: %w", err)
	}

	if isOwner {
		task, err := ts.Task().Get(id)
		if err != nil {
			return fmt.Errorf("task store get: %w", err)
		}

		if taskChanges.Title != nil && *taskChanges.Title != "" {
			task.Title = *taskChanges.Title
		}

		if taskChanges.Description != nil {
			task.Description = taskChanges.Description
		}

		if taskChanges.Status != nil && *taskChanges.Status != "" {
			task.Status = *taskChanges.Status
		}

		err = verifyTask(task)
		if err != nil {
			return fmt.Errorf("verify updated task: %w", err)
		}

		err = ts.WithTx(context.Background(), func(txStore TaskServiceStore) error {
			err = txStore.Task().Update(id, task)
			if err != nil {
				return fmt.Errorf("task store update: %w", err)
			}

			if taskChanges.Assignee != nil {
				err = txStore.Assignee().Create(*taskChanges.Assignee)
				if err != nil {
					return fmt.Errorf("assignee store create: %w", err)
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("WithTx: %w", err)
		}

		return nil
	}

	isAssignee, err := ts.Access().IsAssignee()
	if err != nil {
		return fmt.Errorf("access is assignee: %w", err)
	}

	if isAssignee {
		task, err := ts.Task().Get(id)
		if err != nil {
			return fmt.Errorf("task store get: %w", err)
		}

		if taskChanges.Title != nil && *taskChanges.Title != task.Title {
			return apierr.ErrForbidden
		}
		if taskChanges.Assignee != nil {
			return apierr.ErrForbidden
		}

		if taskChanges.Description != nil {
			task.Description = taskChanges.Description
		}

		if taskChanges.Status != nil && *taskChanges.Status != "" {
			task.Status = *taskChanges.Status
		}

		err = verifyTask(task)
		if err != nil {
			return fmt.Errorf("verify updated task: %w", err)
		}

		err = ts.WithTx(context.Background(), func(txStore TaskServiceStore) error {
			err = txStore.Task().Update(id, task)
			if err != nil {
				return fmt.Errorf("task store update: %w", err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("WithTx: %w", err)
		}

		return nil
	}

	return apierr.ErrForbidden
}
