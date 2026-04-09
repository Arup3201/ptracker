package domain

import "time"

type Assignee struct {
	ProjectID string    `json:"project_id"`
	TaskID    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Avatar    `json:"avatar"`
}
