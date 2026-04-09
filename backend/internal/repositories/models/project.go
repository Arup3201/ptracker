package models

import (
	"time"
)

type Project struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description *string
	Skills      *string
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Tasks        []Task
	Roles        []Membership
	JoinRequests []JoinRequest
	Assignees    []Assignee
	Comments     []Comment
}
