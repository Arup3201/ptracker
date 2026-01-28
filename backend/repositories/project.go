package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/domain"
)

type projectRepo struct {
	DB Execer
}

func NewProjectRepo(db Execer) *projectRepo {
	return &projectRepo{
		DB: db,
	}
}

func (r *projectRepo) Create(ctx context.Context,
	name string,
	description, skills *string,
	owner string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.DB.ExecContext(ctx, "INSERT INTO "+
		"projects(id, name, description, skills, owner, created_at, updated_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $6)",
		id, name, description, skills, owner, now)
	if err != nil {
		return "", fmt.Errorf("db exec context: %w", err)
	}

	return id, nil
}

func (r *projectRepo) All(ctx context.Context, userId string) ([]*domain.ListedProject, error) {
	rows, err := r.DB.QueryContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, r.role, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, p.created_at, p.updated_at "+
			"FROM roles as r "+
			"INNER JOIN projects as p ON r.project_id=p.id "+
			"LEFT JOIN project_summary as ps ON ps.id=p.id "+
			"WHERE r.user_id=($1)",
		userId)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var projects []*domain.ListedProject
	for rows.Next() {
		var p domain.ListedProject
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Role,
			&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
			&p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres get all projects scan: %w", err)
		}
		projects = append(projects, &p)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("rows scan: %w", err)
	}

	return projects, nil
}

func (r *projectRepo) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx,
		"DELETE FROM projects "+
			"WHERE id=($1)",
		id,
	)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
