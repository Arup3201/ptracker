package domain

import "time"

type JoinRequest struct {
	ProjectId string     `json:"project_id"`
	UserId    string     `json:"user_id"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

const (
	JOIN_STATUS_PENDING  = "Pending"
	JOIN_STATUS_ACCEPTED = "Accepted"
	JOIN_STATUS_REJECTED = "Rejected"
)
