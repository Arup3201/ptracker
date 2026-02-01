package domain

import "time"

type Assignee struct {
	ProjectId   string     `json:"project_id"`
	TaskId      string     `json:"task_id"`
	UserId      string     `json:"user_id"`
	Username    string     `json:"username"`
	DisplayName *string    `json:"display_name"`
	Email       string     `json:"email"`
	AvatarURL   *string    `json:"avatar_url"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
