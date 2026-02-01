package domain

import "time"

const (
	ROLE_OWNER  = "Owner"
	ROLE_MEMBER = "Member"
)

type Role struct {
	ProjectId string     `json:"project_id"`
	UserId    string     `json:"user_id"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
