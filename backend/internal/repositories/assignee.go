package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ptracker/internal/apierr"
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

func (r *assigneeRepo) Get(ctx context.Context,
	projectId, taskId, userId string) (bool, error) {

	var tmp int
	err := r.db.QueryRowContext(ctx,
		`SELECT 1 FROM assignees 
		WHERE project_id=($1) AND task_id=($2) AND user_id=($3)`,
		projectId, taskId, userId).Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, apierr.ErrNotFound
	} else if err != nil {
		return false, fmt.Errorf("db query row context: %w", err)
	}

	return true, nil
}

func (r *assigneeRepo) Delete(ctx context.Context,
	projectId, taskId, userId string) error {

	_, err := r.db.ExecContext(ctx,
		"DELETE FROM assignees WHERE project_id=($1) AND task_id=($2) AND user_id=($3)",
		projectId, taskId, userId)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
