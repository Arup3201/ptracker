package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ptracker/models"
)

type ExploredProject struct {
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

type ExploredProjectDetails struct {
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

type ExploreService struct {
	DB     *sql.DB
	UserId string
}

func (ps *ExploreService) GetExploredProjects(page, limit int) ([]ExploredProject, error) {
	var projects = []ExploredProject{}

	rows, err := ps.DB.Query(
		"SELECT "+
			"DISTINCT p.id, p.name, p.description, p.skills, "+
			"CASE "+
			"WHEN r.role='Owner' THEN 'Owner' "+
			"WHEN r.role='Member' THEN 'Member' "+
			"ELSE 'User' "+
			"END AS role, "+
			"p.created_at, p.updated_at "+
			"FROM projects AS p "+
			"CROSS JOIN users AS u "+
			"LEFT JOIN roles AS r ON r.user_id=u.id "+
			"LEFT JOIN join_requests AS jr ON jr.project_id=p.id "+
			"WHERE u.id=($1)",
		ps.UserId)
	if err != nil {
		return projects, fmt.Errorf("postgres get task query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p ExploredProject
		rows.Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Role,
			&p.CreatedAt, &p.UpdatedAt)
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return projects, fmt.Errorf("postgres scan project overview results: %w", err)
	}

	return projects, nil
}

func (ps *ExploreService) RequestToJoinProject(projectId string) error {
	_, err := ps.DB.Exec("INSERT INTO join_requests(user_id, project_id) "+
		"VALUES($1, $2)", ps.UserId, projectId)

	if err != nil {
		return fmt.Errorf("postgres join project: %w", err)
	}

	return nil
}

func (ps *ExploreService) JoinRequests(projectId string) ([]JoinRequest, error) {
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

func (ps *ExploreService) GetProject(projectId string) (*ExploredProjectDetails, error) {
	var p ExploredProjectDetails
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
