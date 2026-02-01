package domain

import "time"

type User struct {
	Id            string     `json:"id"`
	IDPSubject    string     `json:"idp_subject"`
	IDPProvider   string     `json:"idp_provider"`
	Username      string     `json:"username"`
	DisplayName   *string    `json:"display_name"` // nullable
	Email         string     `json:"email"`
	AvatarURL     *string    `json:"avatar_url"` // nullable
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"` // nullable
	LastLoginTime time.Time  `json:"last_login_time"`
}
