package domain

import "time"

const (
	ROLE_OWNER  = "Owner"
	ROLE_MEMBER = "Member"
)

type Role struct {
	ProjectId string
	UserId    string
	Role      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
