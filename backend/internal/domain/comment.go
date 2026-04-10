package domain

import "time"

type Comment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	TaskID    string    `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Avatar    `json:"avatar"`
}
