package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ptracker/apierr"
)

type Assignee struct {
	ProjectId string
	TaskId    string
	UserId    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type AssigneeStore struct {
	DB        *sql.DB
	ProjectId string
	TaskId    string
}

func (s *AssigneeStore) Create(userId string) error {
	_, err := s.DB.Exec(
		"INSERT INTO assignees(project_id, task_id, user_id) "+
			"VALUES($1, $2, $3)",
		s.ProjectId, s.TaskId, userId,
	)

	if err != nil {
		return fmt.Errorf("db insert exec: %w", err)
	}

	return nil
}

func (s *AssigneeStore) Get(userId string) (*Assignee, error) {
	var assignee Assignee
	err := s.DB.QueryRow("SELECT project_id, task_id, user_id, created_at, updated_at "+
		"FROM assignees "+
		"WHERE project_id=($1) AND task_id=($2) AND user_id=($3)",
		s.ProjectId, s.TaskId, userId).
		Scan(&assignee.ProjectId, &assignee.TaskId, &assignee.UserId,
			&assignee.CreatedAt, &assignee.UpdatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apierr.ErrResourceNotFound
		default:
			return nil, fmt.Errorf("db select assignee: %w", err)
		}
	}

	return &assignee, nil
}

func (s *AssigneeStore) Update(userId string) error {
	_, err := s.DB.Exec(
		"UPDATE assignees "+
			"SET user_id=($1) "+
			"WHERE project_id=($2) AND task_id=($3)",
		userId, s.ProjectId, s.TaskId,
	)

	if err != nil {
		return fmt.Errorf("db update exec: %w", err)
	}

	return nil
}

func (s *AssigneeStore) Delete(userId string) error {
	_, err := s.DB.Exec(
		"DELETE FROM assignees "+
			"WHERE project_id=($1) AND task_id=($2) AND user_id=($3)",
		s.ProjectId, s.TaskId, userId,
	)

	if err != nil {
		return fmt.Errorf("db delete exec: %w", err)
	}

	return nil
}
