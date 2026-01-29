package domain

import "time"

type User struct {
	Id            string
	IDPSubject    string
	IDPProvider   string
	Username      string
	DisplayName   *string // nullable
	Email         string
	AvatarURL     *string // nullable
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     *time.Time // nullable
	LastLoginTime time.Time
}
