package domain

import "time"

type User struct {
	ID            string    `json:"id"`
	IDPSubject    string    `json:"-"`
	IDPProvider   string    `json:"-"`
	Username      string    `json:"username"`
	DisplayName   *string   `json:"display_name"` // nullable
	Email         string    `json:"email"`
	AvatarURL     *string   `json:"avatar_url"` // nullable
	IsActive      bool      `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"` // nullable
	LastLoginTime time.Time `json:"-"`
}
