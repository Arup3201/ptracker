package models

import (
	"time"

	"github.com/ptracker/internal/domain"
)

type User struct {
	ID            string `gorm:"primaryKey"`
	IdpSubject    string
	IdpProvider   string
	Username      string
	DisplayName   *string // nullable
	Email         string
	AvatarURL     *string // nullable
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time // nullable
	LastLoginTime time.Time

	Sessions     []Session
	Projects     []Project `gorm:"foreignKey:OwnerID"`
	Memberships  []Membership
	JoinRequests []JoinRequest
	Comments     []Comment
	Assignees    []Assignee
}

func (u User) ToUserDomain() domain.User {
	return domain.User{
		ID:            u.ID,
		IDPSubject:    u.IdpSubject,
		IDPProvider:   u.IdpProvider,
		Username:      u.Username,
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     u.AvatarURL,
		IsActive:      u.IsActive,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		LastLoginTime: u.LastLoginTime,
	}
}
