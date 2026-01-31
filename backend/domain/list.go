package domain

import "time"

type PrivateProjectListed struct {
	*ProjectSummary
	Role string `json:"role"`
}

type PublicProjectListed struct {
	*Project
	Role string `json:"role"`
}

type JoinRequestListed struct {
	ProjectId   string     `json:"project_id"`
	Status      string     `json:"status"`
	UserId      string     `json:"user_id"`
	Username    string     `json:"username"`
	DisplayName *string    `json:"display_name"`
	Email       string     `json:"email"`
	AvatarURL   *string    `json:"avatar_url"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type TaskListed struct {
	*Task
	Assignees []*Assignee `json:"assignees"`
}
