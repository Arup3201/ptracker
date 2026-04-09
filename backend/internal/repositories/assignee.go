package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type assigneeRepo struct {
	db *gorm.DB
}

func NewAssigneeRepo(db *gorm.DB) interfaces.AssigneeRepository {
	return &assigneeRepo{
		db: db,
	}
}

func (r *assigneeRepo) Create(ctx context.Context,
	projectId, taskId, userId string) error {

	var err error

	assignee := models.Assignee{
		ProjectID: projectId,
		TaskID:    taskId,
		UserID:    userId,
	}
	err = gorm.G[models.Assignee](r.db).Create(ctx, &assignee)
	if err != nil {
		return fmt.Errorf("gorm create: %w", err)
	}

	return nil
}

func (r *assigneeRepo) Is(ctx context.Context,
	projectId, taskId, userId string) (bool, error) {

	assignees, err := gorm.G[models.Assignee](r.db).Where(
		"project_id = ? AND task_id = ? AND user_id = ?",
		projectId, taskId, userId).Find(ctx)
	if err != nil {
		return false, fmt.Errorf("gorm query: %w", err)
	}

	if len(assignees) > 0 {
		return true, nil
	}

	return false, nil
}

func (r *assigneeRepo) Delete(ctx context.Context,
	projectId, taskId, userId string) error {

	_, err := gorm.G[models.Assignee](r.db).Where(
		"project_id = ? AND task_id = ? AND user_id = ?",
		projectId, taskId, userId).Delete(ctx)
	if err != nil {
		return fmt.Errorf("gorm query delete: %w", err)
	}

	return nil
}
