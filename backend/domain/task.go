package domain

import "time"

type Task struct {
	Id          string
	ProjectId   string
	Title       string
	Description *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
