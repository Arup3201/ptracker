package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type PublicRepo struct {
	db interfaces.Execer
}

func NewPublicRepo(db interfaces.Execer) interfaces.PublicRepository {
	return &PublicRepo{
		db: db,
	}
}

func (r *PublicRepo) Get(ctx context.Context, projectId string) (*domain.PublicProjectSummary, error) {
	var p = domain.PublicProjectSummary{
		ProjectSummary: &domain.ProjectSummary{
			Project: &domain.Project{},
		},
		Owner: &domain.Member{},
	}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, p.owner, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, "+
			"CASE "+
			"WHEN jr.status='Pending' THEN 'Pending' "+
			"WHEN jr.status='Accepted' THEN 'Accepted' "+
			"ELSE 'Not Requested' "+
			"END AS join_status, "+
			"p.created_at, p.updated_at "+
			"FROM projects as p "+
			"LEFT JOIN project_summary as ps ON p.id=ps.id "+
			"LEFT JOIN join_requests AS jr ON p.id=jr.project_id "+
			"WHERE p.id=($1)",
		projectId).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner.Id,
		&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
		&p.JoinStatus, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("db query row context: %w", err)
	}

	return &p, nil
}
