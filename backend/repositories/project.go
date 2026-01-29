package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
)

type projectRepo struct {
	db Execer
}

func NewProjectRepo(db Execer) stores.ProjectRepository {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) Create(ctx context.Context,
	name string,
	description, skills *string,
	owner string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.db.ExecContext(ctx, "INSERT INTO "+
		"projects(id, name, description, skills, owner, created_at, updated_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $6)",
		id, name, description, skills, owner, now)
	if err != nil {
		return "", fmt.Errorf("db exec context: %w", err)
	}

	return id, nil
}

func (r *projectRepo) Get(ctx context.Context, projectId string) (*domain.ProjectSummary, error) {
	var p domain.ProjectSummary
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, "+
			"p.created_at, p.updated_at "+
			"FROM projects as p "+
			"LEFT JOIN project_summary as ps ON p.id=ps.id "+
			"WHERE p.id=($1)",
		projectId).Scan(&p.Id, &p.Name, &p.Description, &p.Skills,
		&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("postgres query project details: %w", err)
	}

	return &p, nil
}

func (r *projectRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM projects "+
			"WHERE id=($1)",
		id,
	)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
