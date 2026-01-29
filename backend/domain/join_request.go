package domain

import "time"

type JoinRequest struct {
	ProjectId string
	UserId    string
	Status    string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

const (
	JOIN_STATUS_PENDING  = "Pending"
	JOIN_STATUS_ACCEPTED = "Accepted"
	JOIN_STATUS_REJECTED = "Rejected"
)
