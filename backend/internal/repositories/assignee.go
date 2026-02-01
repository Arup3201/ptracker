package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/ptracker/internal/interfaces"
)

type assigneeRepo struct {
	db interfaces.Execer
}

func NewAssigneeRepo(db interfaces.Execer) interfaces.AssigneeRepository {
	return &assigneeRepo{
		db: db,
	}
}

func (r *assigneeRepo) Create(ctx context.Context,
	projectId, taskId, userId string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO 
		assignees(project_id, task_id, user_id, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$4)`,
		projectId, taskId, userId, now)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
