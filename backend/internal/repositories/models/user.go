package models

import (
	"time"

	"github.com/ptracker/internal/domain"
)

type User struct {
	ID            string `gorm:"primaryKey"`
	IdpSubject    string `gorm:"index:idx_user_idp,unique"`
	IdpProvider   string `gorm:"index:idx_user_idp,unique"`
	Username      string
	DisplayName   *string // nullable
	Email         string
	AvatarURL     *string // nullable
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time // nullable
	LastLoginTime time.Time

	Sessions     []Session     `gorm:"constraint:OnDelete:CASCADE"`
	Projects     []Project     `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Memberships  []Membership  `gorm:"constraint:OnDelete:CASCADE"`
	JoinRequests []JoinRequest `gorm:"constraint:OnDelete:CASCADE"`
	Comments     []Comment     `gorm:"constraint:OnDelete:CASCADE"`
	Assignees    []Assignee    `gorm:"constraint:OnDelete:CASCADE"`
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
