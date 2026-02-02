package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type TaskRepo struct {
	db interfaces.Execer
}

func NewTaskRepo(db interfaces.Execer) interfaces.TaskRepository {
	return &TaskRepo{
		db: db,
	}
}

func (r *TaskRepo) Create(ctx context.Context,
	projectId string,
	title, description, status string) (string, error) {
	tId := uuid.NewString()
	now := time.Now()

	_, err := r.db.ExecContext(ctx,
		"INSERT INTO "+
			"tasks(id, project_id, title, description, status, created_at, updated_at) "+
			"VALUES($1, $2, $3, $4, $5, $6, $6)",
		tId, projectId, title, description, status, now)
	if err != nil {
		return "", fmt.Errorf("insert task: %w", err)
	}

	return tId, nil
}

func (r *TaskRepo) Get(ctx context.Context, id string) (*domain.Task, error) {
	var pt domain.Task
	err := r.db.QueryRowContext(ctx,
		`SELECT t.id, t.title, t.description, 
			t.status, t.created_at, t.updated_at 
			FROM tasks as t 
			WHERE t.id=($1) AND deleted_at IS NULL`,
		id).
		Scan(&pt.Id, &pt.Title, &pt.Description, &pt.Status, &pt.CreatedAt, &pt.UpdatedAt)
	if err != nil {
		return &pt, fmt.Errorf("db query row context: %w", err)
	}

	return &pt, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {

	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks 
		SET 
		title=($2), description=($3), status=($4), updated_at=($5) 
		WHERE id=($1)`,
		task.Id,
		task.Title, task.Description, task.Status, now)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
