package domain

import "time"

type Task struct {
	Id          string     `json:"id"`
	ProjectId   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type TaskDetail struct {
	*Task
	Assignees []*Assignee `json:"assignees"`
}

const (
	TASK_STATUS_UNASSIGNED = "Unassigned"
	TASK_STATUS_ONGOING    = "Ongoing"
	TASK_STATUS_COMPLETED  = "Completed"
	TASK_STATUS_ABANDONED  = "Abandoned"
)
