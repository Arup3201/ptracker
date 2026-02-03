package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type ListRepo struct {
	db interfaces.Execer
}

func NewListRepo(db interfaces.Execer) interfaces.ListRepository {
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
		rows.Scan(&m.UserId, &m.Username, &m.DisplayName, &m.Email,
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

func (r *ListRepo) Comments(ctx context.Context,
	projectId, taskId string) ([]*domain.Comment, error) {

	rows, err := r.db.QueryContext(ctx,
		`SELECT 
			c.id, c.project_id, c.task_id, c.content, 
			u.id, u.username, u.display_name, u.email, 
			u.avatar_url, u.is_active, 
			r.created_at, r.updated_at,
			c.created_at, c.updated_at 
			FROM comments AS c 
			INNER JOIN users AS u ON c.user_id=u.id 
			INNER JOIN roles AS r ON c.project_id=r.project_id AND c.user_id=r.user_id 
			WHERE c.project_id=($1) AND c.task_id=($2)`,
		projectId, taskId,
	)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var comments = []*domain.Comment{}
	for rows.Next() {
		var (
			comment domain.Comment
			user    domain.Member
		)

		err := rows.Scan(&comment.Id, &comment.ProjectId, &comment.TaskId,
			&comment.Content,
			&user.UserId, &user.Username, &user.DisplayName, &user.Email,
			&user.AvatarURL, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
			&comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		comment.User = &user
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return comments, nil
}

func (r *ListRepo) RecentlyCreatedProjects(ctx context.Context,
	userId string, n int) ([]*domain.RecentProjectListed, error) {

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, "+
			"p.created_at, p.updated_at "+
			"FROM projects AS p "+
			"LEFT JOIN project_summary as ps ON ps.id=p.id "+
			"WHERE p.owner=($1) "+
			"ORDER BY p.created_at DESC "+
			"LIMIT ($2)",
		userId, n)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var projects = []*domain.RecentProjectListed{}
	for rows.Next() {
		var (
			project                                   domain.Project
			unassigned, ongoing, completed, abandoned int
		)

		err := rows.Scan(&project.Id, &project.Name, &project.Description, &project.Skills,
			&unassigned, &ongoing, &completed, &abandoned,
			&project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		projects = append(projects, &domain.RecentProjectListed{
			ProjectSummary: &domain.ProjectSummary{
				Project:         &project,
				UnassignedTasks: unassigned,
				OngoingTasks:    ongoing,
				CompletedTasks:  completed,
				AbandonedTasks:  abandoned,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("rows: %w", err)
	}

	return projects, nil
}

func (r *ListRepo) RecentlyJoinedProjects(ctx context.Context,
	userId string,
	n int) ([]*domain.RecentProjectListed, error) {

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+
			"p.id, p.name, p.description, p.skills, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, "+
			"p.created_at, p.updated_at "+
			"FROM roles AS r "+
			"INNER JOIN projects AS p ON p.id=r.project_id "+
			"LEFT JOIN project_summary as ps ON p.id=ps.id "+
			"WHERE r.user_id=($1) AND r.role=($2) "+
			"ORDER BY r.created_at DESC "+
			"LIMIT ($3)",
		userId, domain.ROLE_MEMBER, n)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()
	var projects = []*domain.RecentProjectListed{}
	for rows.Next() {
		var (
			project                                   domain.Project
			unassigned, ongoing, completed, abandoned int
		)

		err := rows.Scan(&project.Id, &project.Name, &project.Description, &project.Skills,
			&unassigned, &ongoing, &completed, &abandoned,
			&project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}
		projects = append(projects, &domain.RecentProjectListed{
			ProjectSummary: &domain.ProjectSummary{
				Project:         &project,
				UnassignedTasks: unassigned,
				OngoingTasks:    ongoing,
				CompletedTasks:  completed,
				AbandonedTasks:  abandoned,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("rows: %w", err)
	}

	return projects, nil
}

func (r *ListRepo) RecentlyAssignedTasks(ctx context.Context,
	userId string,
	n int) ([]*domain.RecentTaskListed, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			t.id, t.project_id, t.title, t.status, t.created_at, t.updated_at, 
			p.name
			FROM tasks AS t 
			INNER JOIN projects AS p ON t.project_id=p.id
			INNER JOIN assignees AS a ON a.task_id=t.id 
			WHERE a.user_id=($1)  
			ORDER BY a.created_at DESC 
			LIMIT ($2)`, userId, n)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var tasks = []*domain.RecentTaskListed{}
	for rows.Next() {
		var (
			task        domain.Task
			projectName string
		)
		err := rows.Scan(&task.Id, &task.ProjectId, &task.Title,
			&task.Status, &task.CreatedAt, &task.UpdatedAt,
			&projectName)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		tasks = append(tasks, &domain.RecentTaskListed{
			Task:        &task,
			ProjectName: projectName,
		})
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (r *ListRepo) RecentlyUnassignedTasks(ctx context.Context,
	userId string,
	n int) ([]*domain.RecentTaskListed, error) {

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			t.id, t.project_id, t.title, t.status, t.created_at, t.updated_at, 
			p.name 
			FROM tasks AS t 
			INNER JOIN projects AS p ON t.project_id=p.id 
			WHERE p.owner=($1) AND t.status=($2) 
			ORDER BY t.created_at DESC 
			LIMIT ($3)`, userId, domain.TASK_STATUS_UNASSIGNED, n)
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}
	defer rows.Close()

	var tasks = []*domain.RecentTaskListed{}
	for rows.Next() {
		var (
			task        domain.Task
			projectName string
		)
		err := rows.Scan(&task.Id, &task.ProjectId, &task.Title,
			&task.Status, &task.CreatedAt, &task.UpdatedAt, &projectName)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		tasks = append(tasks, &domain.RecentTaskListed{
			Task:        &task,
			ProjectName: projectName,
		})
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}
