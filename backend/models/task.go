package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProjectTask struct {
	Id          string
	ProjectId   string
	Title       string
	Description *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type TaskStore struct {
	DB        *sql.DB
	UserId    string
	ProjectId string
}

func (ts *TaskStore) Create(title, description, status string) (string, error) {
	tId := uuid.NewString()
	_, err := ts.DB.Exec("INSERT INTO "+
		"tasks(id, project_id, title, description, status) "+
		"VALUES($1, $2, $3, $4, $5)",
		tId, ts.ProjectId, title, description, status)
	if err != nil {
		return "", fmt.Errorf("insert task: %w", err)
	}

	return tId, nil
}

func (ts *TaskStore) All() ([]ProjectTask, error) {
	rows, err := ts.DB.Query("SELECT t.id, t.title, "+
		"t.status, t.created_at, t.updated_at "+
		"FROM tasks as t "+
		"WHERE t.project_id=($1) AND deleted_at IS NULL", ts.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("postgres get all projects query: %w", err)
	}
	defer rows.Close()

	var tasks []ProjectTask
	for rows.Next() {
		var t ProjectTask
		err := rows.Scan(&t.Id, &t.Title,
			&t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres get all projects scan: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (ts *TaskStore) Count() (int, error) {
	var cnt int
	err := ts.DB.QueryRow("SELECT COUNT(t.id) "+
		"FROM projects as p "+
		"INNER JOIN tasks as t ON p.id=t.project_id "+
		"WHERE p.id=($1)", ts.ProjectId).Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("postgres get task count: %w", err)
	}

	return cnt, nil
}

func (ts *TaskStore) Get(id string) (ProjectTask, error) {
	var pt ProjectTask
	err := ts.DB.QueryRow("SELECT t.id, t.title, t.description, "+
		"t.status, t.created_at, t.updated_at "+
		"FROM tasks as t "+
		"WHERE t.id=($1) AND t.project_id=($2) AND deleted_at IS NULL", id, ts.ProjectId).
		Scan(&pt.Id, &pt.Title, &pt.Description, &pt.Status, &pt.CreatedAt, &pt.UpdatedAt)
	if err != nil {
		return pt, fmt.Errorf("postgres get task query: %w", err)
	}

	return pt, nil
}

func (ts *TaskStore) Update(id string, changes *ProjectTask) error {
	_, err := ts.DB.Exec("UPDATE tasks "+
		"SET title=($1), description=($2), status=($3), updated_at=CURRENT_TIMESTAMP "+
		"WHERE id=($4)", changes.Title, changes.Description, changes.Status, id)
	if err != nil {
		return fmt.Errorf("db exec: %w", err)
	}

	return nil
}
