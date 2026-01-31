package domain

import "time"

type Member struct {
	Id          string     `json:"id"`
	Username    string     `json:"username"`
	DisplayName *string    `json:"display_name"`
	Email       string     `json:"email"`
	AvatarURL   *string    `json:"avatar_url"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
