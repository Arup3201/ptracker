package domain

import "time"

type Comment struct {
	Id        string     `json:"id"`
	ProjectId string     `json:"project_id"`
	TaskId    string     `json:"task_id"`
	User      *Member    `json:"user"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
