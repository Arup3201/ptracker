package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ptracker/apierr"
	"github.com/ptracker/models"
)

type Project struct {
	Id          string
	Name        string
	Description *string
	Skills      *string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type JoinRequest struct {
	ProjectId string
	User      models.User
	Status    string
	CreatedAt string
}

type ProjectDetails struct {
	Id              string
	Name            string
	Description     *string
	Skills          *string
	Owner           string
	JoinStatus      string
	UnassignedTasks int
	OngoingTasks    int
	CompletedTasks  int
	AbandonedTasks  int
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type Member struct {
	ProjectId   string
	UserId      string
	Username    string
	DisplayName string
	Email       string
	AvaterURL   *string
	IsActive    bool
	Role        string
	JoinedAt    time.Time
}

type ProjectService struct {
	DB     *sql.DB
	UserId string
}

func (ps *ProjectService) List(page, limit int) ([]Project, error) {
	var projects = []Project{}

	rows, err := ps.DB.Query(
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
		ps.UserId)
	if err != nil {
		return projects, fmt.Errorf("postgres get task query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		rows.Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Role,
			&p.CreatedAt, &p.UpdatedAt)
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("postgres scan project overview results: %w", err)
	}

	return projects, nil
}

func (ps *ProjectService) Join(projectId string) error {
	_, err := ps.DB.Exec("INSERT INTO join_requests(user_id, project_id) "+
		"VALUES($1, $2)", ps.UserId, projectId)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return apierr.ErrDuplicate
		}
		return fmt.Errorf("postgres join project: %w", err)
	}

	return nil
}

func (ps *ProjectService) JoinRequests(projectId string) ([]JoinRequest, error) {
	rows, err := ps.DB.Query("SELECT r.project_id, r.status, u.id, u.username, "+
		"u.display_name, u.email, u.is_active, r.created_at "+
		"FROM join_requests as r "+
		"INNER JOIN users as u ON u.id=r.user_id "+
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

	results := []JoinRequest{}
	for rows.Next() {
		var r JoinRequest
		rows.Scan(&r.ProjectId, &r.Status, &r.User.Id,
			&r.User.Username, &r.User.DisplayName, &r.User.Email,
			&r.User.IsActive, &r.CreatedAt)

		results = append(results, r)
	}

	return results, nil
}

func (ps *ProjectService) Get(projectId string) (*ProjectDetails, error) {
	var p ProjectDetails
	err := ps.DB.QueryRow(
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
		projectId).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner,
		&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
		&p.JoinStatus, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("postgres query project details: %w", err)
	}

	return &p, nil
}

func (ps *ProjectService) UpdateJoinRequestStatus(projectId, userId, joinStatus string) error {
	if joinStatus == "Pending" {
		return apierr.ErrInvalidValue
	}

	ctx := context.Background()
	tx, err := ps.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("database transaction begin: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE join_requests "+
		"SET status=($1) "+
		"WHERE project_id=($2) AND user_id=($3)", joinStatus, projectId, userId)

	if err != nil {
		if strings.Contains(err.Error(), "invalid input value") {
			return apierr.ErrInvalidValue
		}
		return fmt.Errorf("service update join request query: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO roles(project_id, user_id, role) "+
		"VALUES($1, $2, $3)", projectId, userId, models.ROLE_MEMBER)
	if err != nil {
		return fmt.Errorf("service insert member role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit: %w", err)
	}

	return nil
}

func (ps *ProjectService) Members(projectId string) ([]Member, error) {
	var members []Member

	rows, err := ps.DB.Query("SELECT "+
		"r.project_id, u.id, u.username, u.display_name, u.email, "+
		"u.avatar_url, u.is_active, r.role, r.created_at "+
		"FROM roles AS r "+
		"INNER JOIN users AS u ON r.user_id=u.id "+
		"WHERE r.project_id=($1) AND r.role!=($2)", projectId, models.ROLE_OWNER)
	if err != nil {
		return nil, fmt.Errorf("service get members query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m Member
		rows.Scan(&m.ProjectId, &m.UserId, &m.Username, &m.DisplayName, &m.Email,
			&m.AvaterURL, &m.IsActive, &m.Role, &m.JoinedAt)

		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("service rows scan: %w", err)
	}

	return members, nil
}
