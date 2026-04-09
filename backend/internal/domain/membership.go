package domain

import "time"

const (
	ROLE_OWNER  = "Owner"
	ROLE_MEMBER = "Member"
)

type Membership struct {
	ProjectID string    `json:"project_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Avatar    `json:"avatar"`
}
