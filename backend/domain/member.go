package domain

import "time"

type Member struct {
	Id          string
	Username    string
	DisplayName *string
	Email       string
	AvaterURL   *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
