package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
)

type ListRepo struct {
	db Execer
}

func NewListRepo(db Execer) interfaces.ListRepository {
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
		var p = domain.PrivateProjectListed{
			ProjectSummary: &domain.ProjectSummary{
				Project: &domain.Project{},
			},
		}

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

func (r *ListRepo) Members(ctx context.Context,
	projectId string) ([]*domain.Member, error) {
	var members []*domain.Member

	rows, err := r.db.QueryContext(ctx,
		"SELECT "+
			"u.id, u.username, u.display_name, u.email, "+
			"u.avatar_url, u.is_active, r.created_at, r.updated_at "+
			"FROM roles AS r "+
			"INNER JOIN users AS u ON r.user_id=u.id "+
			"WHERE r.project_id=($1) AND r.role!=($2)", projectId, domain.ROLE_OWNER)
	if err != nil {
		return nil, fmt.Errorf("service get members query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m domain.Member
		rows.Scan(&m.Id, &m.Username, &m.DisplayName, &m.Email,
			&m.AvatarURL, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)

		members = append(members, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("service rows scan: %w", err)
	}

	return members, nil
}

func (r *ListRepo) PublicProjects(ctx context.Context, userId string) ([]*domain.PublicProjectListed, error) {
	var projects = []*domain.PublicProjectListed{}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, "+
			"CASE "+
			"WHEN p.owner=($1) THEN 'Owner' "+
			"WHEN r.user_id=($1) THEN 'Member' "+
			"ELSE 'User' "+
			"END AS role, "+
			"p.created_at, p.updated_at "+
			"FROM projects AS p "+
			"LEFT JOIN roles AS r ON p.id=r.project_id AND r.user_id=($1)",
		userId)
	if err != nil {
		return projects, fmt.Errorf("postgres get task query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p = domain.PublicProjectListed{
			Project: &domain.Project{},
		}
		rows.Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Role,
			&p.CreatedAt, &p.UpdatedAt)
		projects = append(projects, &p)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("postgres scan project overview results: %w", err)
	}

	return projects, nil
}

func (r *ListRepo) JoinRequests(ctx context.Context, projectId string) ([]*domain.JoinRequestListed, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT r.project_id, r.status, jr.created_at, jr.updated_at, u.id, u.username, "+
			"u.display_name, u.email, u.is_active, r.created_at, r.updated_at "+
			"FROM join_requests as jr "+
			"INNER JOIN users as u ON u.id=jr.user_id "+
			"INNER JOIN roles as r ON r.user_id=jr.user_id AND r.project_id=jr.project_id "+
			"WHERE r.project_id=($1)",
		projectId,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres get join requests: %w", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres get join requests: %w", err)
	}

	results := []*domain.JoinRequestListed{}
	for rows.Next() {
		var r = domain.JoinRequestListed{
			JoinRequest: &domain.JoinRequest{},
			Member:      &domain.Member{},
		}
		rows.Scan(&r.JoinRequest.ProjectId, &r.JoinRequest.Status,
			&r.JoinRequest.CreatedAt, &r.JoinRequest.UpdatedAt,
			&r.Member.Id, &r.Member.Username, &r.Member.DisplayName,
			&r.Member.Email, &r.Member.IsActive,
			&r.Member.CreatedAt, &r.Member.UpdatedAt)

		results = append(results, &r)
	}

	return results, nil
}
