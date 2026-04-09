package models

import (
	"time"
)

type Task struct {
	ID          string `gorm:"primaryKey"`
	ProjectID   string `gorm:"index:idx_task_project"`
	Title       string
	Description *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Assignees []Assignee
	Comments  []Comment
}
