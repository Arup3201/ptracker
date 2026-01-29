package domain

import "time"

type Member struct {
	Id          string
	Username    string
	DisplayName *string
	Email       string
	AvatarURL   *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
