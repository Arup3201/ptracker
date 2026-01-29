package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/domain"
)

type ListRepo struct {
	db Execer
}

func NewListRepo(db Execer) *ListRepo {
	return &ListRepo{
		db: db,
	}
}

func (r *ListRepo) PrivateProjects(ctx context.Context, userId string) ([]*domain.PrivateProjectListed, error) {
	rows, err := r.db.QueryContext(
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

	var projects []*domain.PrivateProjectListed
	for rows.Next() {
		var p domain.PrivateProjectListed
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
