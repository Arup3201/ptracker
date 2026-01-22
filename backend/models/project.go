package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/apierr"
)

type Project struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	Skills          *string `json:"skills"`
	Owner           string
	Role            string     `json:"role"`
	UnassignedTasks int        `json:"unassigned_tasks"`
	OngoingTasks    int        `json:"ongoing_tasks"`
	CompletedTasks  int        `json:"completed_tasks"`
	AbandonedTasks  int        `json:"abandoned_tasks"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

const (
	ROLE_OWNER  = "Owner"
	ROLE_MEMBER = "Member"
)

const (
	TASK_STATUS_UNASSIGNED = "Unassigned"
	TASK_STATUS_ONGOING    = "Ongoing"
	TASK_STATUS_COMPLETED  = "Completed"
	TASK_STATUS_ABANDONED  = "Abandoned"
)

var (
	TASK_STATUS = []string{TASK_STATUS_UNASSIGNED, TASK_STATUS_ONGOING,
		TASK_STATUS_COMPLETED, TASK_STATUS_ABANDONED}
)

type ProjectStore struct {
	DB     *sql.DB
	UserId string
}

// Create a project in the database, returns the ID of the created project.
func (ps *ProjectStore) Create(name, description, skills string) (string, error) {
	ctx := context.Background()
	tx, err := ps.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("transaction begin: %w", err)
	}
	defer tx.Rollback()

	// insert project row
	pid := uuid.NewString()
	_, err = tx.ExecContext(ctx, "INSERT INTO "+
		"projects(id, name, description, skills, owner) "+
		"VALUES($1, $2, $3, $4, $5)",
		pid, name, description, skills, ps.UserId)
	if err != nil {
		return "", fmt.Errorf("insert project: %w", err)
	}

	// insert role as "Owner"
	_, err = tx.ExecContext(ctx, "INSERT INTO "+
		"roles(user_id, project_id, role) "+
		"VALUES($1, $2, $3)",
		ps.UserId, pid, ROLE_OWNER)
	if err != nil {
		return "", fmt.Errorf("insert role: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("transaction commit: %w", err)
	}

	return pid, nil
}

func (ps *ProjectStore) All(page, limit int) ([]Project, error) {
	rows, err := ps.DB.Query(
		"SELECT "+
			"p.id, p.name, p.description, p.skills, r.role, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, p.created_at, p.updated_at "+
			"FROM roles as r "+
			"INNER JOIN projects as p ON r.project_id=p.id "+
			"LEFT JOIN project_summary as ps ON ps.id=p.id "+
			"WHERE r.user_id=($1)",
		ps.UserId)
	if err != nil {
		return nil, fmt.Errorf("postgres get all projects query: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Role,
			&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
			&p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres get all projects scan: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return projects, err
	}

	return projects, nil
}

func (ps *ProjectStore) Count() (int, error) {
	var cnt int
	err := ps.DB.QueryRow("SELECT COUNT(p.id) "+
		"FROM projects as p "+
		"LEFT JOIN roles as r ON p.id=r.project_id "+
		"WHERE r.user_id=($1)"+
		"GROUP BY p.id", ps.UserId).Scan(&cnt)
	if err == sql.ErrNoRows {
		return 0, apierr.ErrResourceNotFound
	} else if err != nil {
		return 0, fmt.Errorf("postgres get active session: %w", err)
	}

	return cnt, nil
}

func (ps *ProjectStore) Get(projectId string) (*Project, error) {
	var p Project
	err := ps.DB.QueryRow(
		"SELECT "+
			"p.id, p.name, p.description, p.skills, p.owner, "+
			"ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, "+
			"p.created_at, p.updated_at "+
			"FROM projects as p "+
			"LEFT JOIN project_summary as ps ON p.id=ps.id "+
			"WHERE p.id=($1)",
		projectId).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner,
		&p.UnassignedTasks, &p.OngoingTasks, &p.CompletedTasks, &p.AbandonedTasks,
		&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("postgres query project details: %w", err)
	}

	return &p, nil
}

func (ps *ProjectStore) CountMembers(projectId string) (int, error) {
	var count int
	err := ps.DB.QueryRow(
		"SELECT "+
			"COUNT(user_id) "+
			"FROM roles "+
			"WHERE user_id=($1) AND project_id=($2) AND role!=($3)",
		ps.UserId, projectId, ROLE_OWNER,
	).Scan(count)
	if err != nil {
		return -1, fmt.Errorf("postgres query total members: %w", err)
	}

	return count, nil
}
