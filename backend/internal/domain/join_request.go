package domain

import "time"

const (
	JOIN_STATUS_PENDING  = "Pending"
	JOIN_STATUS_ACCEPTED = "Accepted"
	JOIN_STATUS_REJECTED = "Rejected"
)

type JoinRequest struct {
	ProjectID string    `json:"project_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Avatar    `json:"avatar"`
}
