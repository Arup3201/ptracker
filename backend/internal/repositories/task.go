package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) interfaces.TaskRepository {
	return &TaskRepo{
		db: db,
	}
}

func (r *TaskRepo) Create(ctx context.Context,
	projectId string,
	title, description, status string) (string, error) {

	id := uuid.NewString()
	task := models.Task{
		ID:          id,
		ProjectID:   projectId,
		Title:       title,
		Description: &description,
		Status: models.TaskStatus{
			String: status,
		},
	}

	err := gorm.G[models.Task](r.db).Create(ctx, &task)
	if err != nil {
		return "", fmt.Errorf("gorm create: %w", err)
	}

	return task.ID, nil
}

func (r *TaskRepo) Get(ctx context.Context, id string) (domain.ProjectTaskItem, error) {

	query := `SELECT 
			t.id, t.title, t.status, t.created_at, t.updated_at, 
			COALESCE(
				json_agg(
				json_build_object(
					'project_id', a.project_id,
					'task_id', a.task_id,
					'assignee_id', a.user_id,
					'assignee_username', u.username,
					'assignee_display_name', u.display_name,
					'assignee_email', u.email,
					'assignee_avatar_url', u.avatar_url
				)
				) FILTER (WHERE a.user_id IS NOT NULL),
				'[]'
			) AS assignees 
			FROM tasks AS t 
			LEFT JOIN assignees AS a ON a.task_id=t.id 
			LEFT JOIN users AS u ON u.id=a.user_id 
			WHERE t.id = ? 
			GROUP BY t.id, t.title, t.status, t.created_at, t.updated_at`

	task, err := gorm.G[models.ProjectTaskItemRow](r.db).Raw(query, id).First(ctx)
	if err != nil {
		return task.ToProjectTaskItemDomain(), fmt.Errorf("gorm query task: %w", err)
	}

	return task.ToProjectTaskItemDomain(), nil
}

func (r *TaskRepo) Update(ctx context.Context, id string,
	title, description, status *string) error {

	task, err := gorm.G[models.Task](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return fmt.Errorf("gorm query task: %w", err)
	}

	if title != nil {
		task.Title = *title
	}

	if description != nil {
		task.Description = description
	}

	if status != nil {
		task.Status = models.TaskStatus{
			String: *status,
		}
	}

	err = r.db.Save(&task).Error
	if err != nil {
		return fmt.Errorf("gorm db save: %w", err)
	}

	return nil
}
