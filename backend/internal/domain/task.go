package domain

import "time"

const (
	TASK_STATUS_UNASSIGNED = "Unassigned"
	TASK_STATUS_ONGOING    = "Ongoing"
	TASK_STATUS_COMPLETED  = "Completed"
	TASK_STATUS_ABANDONED  = "Abandoned"
)

type ProjectTaskItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ProjectID string     `json:"project_id"`
	Assignees []Assignee `json:"assignees"`
}

type DashboardTaskItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
}
