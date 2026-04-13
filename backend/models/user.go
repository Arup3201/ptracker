package models

import "time"

type User struct {
	ID            string `gorm:"primaryKey"`
	Username      string
	DisplayName   *string // nullable
	Email         string
	AvatarURL     *string // nullable
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLoginTime time.Time

	Projects     []Project     `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Members      []Member      `gorm:"constraint:OnDelete:CASCADE"`
	JoinRequests []JoinRequest `gorm:"constraint:OnDelete:CASCADE"`
	Comments     []Comment     `gorm:"constraint:OnDelete:CASCADE"`
	Assignees    []Assignee    `gorm:"constraint:OnDelete:CASCADE"`
}
