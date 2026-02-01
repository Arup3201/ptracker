package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
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

	var projects = []*domain.PrivateProjectListed{}
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

func (r *ListRepo) Tasks(ctx context.Context,
	projectId string) ([]*domain.TaskListed, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			t.id, t.title, t.status, t.created_at, t.updated_at, 
			COALESCE(
				json_agg(
				json_build_object(
					'project_id', a.project_id,
					'task_id', a.task_id,
					'user_id', a.user_id,
					'username', u.username,
					'display_name', u.display_name,
					'email', u.email,
					'avatar_url', u.avatar_url,
					'is_active', u.is_active,
					'created_at', a.created_at,
					'updated_at', a.updated_at
				)
				) FILTER (WHERE a.user_id IS NOT NULL),
				'[]'
			) AS assignees 
			FROM tasks AS t 
			LEFT JOIN assignees AS a ON a.task_id=t.id 
			LEFT JOIN users AS u ON u.id=a.user_id 
			WHERE t.project_id=($1) AND t.deleted_at IS NULL
			GROUP BY t.id`, projectId)
	if err != nil {
		return nil, fmt.Errorf("postgres get all projects query: %w", err)
	}
	defer rows.Close()

	var tasks = []*domain.TaskListed{}
	for rows.Next() {
		var (
			task      domain.Task
			assignees []byte
		)
		err := rows.Scan(&task.Id, &task.Title,
			&task.Status, &task.CreatedAt, &task.UpdatedAt,
			&assignees)
		if err != nil {
			return nil, fmt.Errorf("postgres get all projects scan: %w", err)
		}

		var list []*domain.Assignee
		if err := json.Unmarshal(assignees, &list); err != nil {
			return nil, err
		}

		tasks = append(tasks, &domain.TaskListed{
			Task:      &task,
			Assignees: list,
		})
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (r *ListRepo) Assignees(ctx context.Context,
	taskId string) ([]*domain.Assignee, error) {

	rows, err := r.db.QueryContext(ctx,
		"SELECT "+
			"a.project_id, a.task_id, u.id, u.username, u.display_name, "+
			"u.avatar_url, u.is_active, a.created_at, a.updated_at "+
			"FROM assignees AS a "+
			"INNER JOIN users AS u ON a.user_id=u.id "+
			"WHERE a.task_id=($1)",
		taskId,
	)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var assignees = []*domain.Assignee{}
	for rows.Next() {
		var a domain.Assignee

		err := rows.Scan(&a.ProjectId, &a.TaskId, &a.UserId, &a.Username,
			&a.DisplayName, &a.AvatarURL, &a.IsActive,
			&a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		assignees = append(assignees, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return assignees, nil
}

func (r *ListRepo) Members(ctx context.Context,
	projectId string) ([]*domain.Member, error) {

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

	var members = []*domain.Member{}
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
		return nil, fmt.Errorf("postgres get task query: %w", err)
	}
	defer rows.Close()

	var projects = []*domain.PublicProjectListed{}
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
		"SELECT jr.project_id, jr.status, u.id, u.username, "+
			"u.display_name, u.email, u.avatar_url, u.is_active, "+
			"jr.created_at, jr.updated_at "+
			"FROM join_requests as jr "+
			"INNER JOIN users as u ON u.id=jr.user_id "+
			"WHERE jr.project_id=($1) ",
		projectId,
	)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	results := []*domain.JoinRequestListed{}
	for rows.Next() {
		var r = domain.JoinRequestListed{}
		rows.Scan(&r.ProjectId, &r.Status,
			&r.UserId, &r.Username, &r.DisplayName,
			&r.Email, &r.AvatarURL, &r.IsActive,
			&r.CreatedAt, &r.UpdatedAt)

		results = append(results, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return results, nil
}
