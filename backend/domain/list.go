package domain

import "time"

type PrivateProjectListed struct {
	*ProjectSummary
	Role string
}

type PublicProjectListed struct {
	*Project
	Role string
}

type JoinRequestListed struct {
	ProjectId   string
	Status      string
	UserId      string
	Username    string
	DisplayName *string
	Email       string
	AvatarURL   *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
