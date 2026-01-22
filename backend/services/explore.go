package services

import (
	"database/sql"
	"fmt"
	"time"
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
