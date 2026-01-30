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

const (
	TASK_STATUS_UNASSIGNED = "Unassigned"
	TASK_STATUS_ONGOING    = "Ongoing"
	TASK_STATUS_COMPLETED  = "Completed"
	TASK_STATUS_ABANDONED  = "Abandoned"
)
