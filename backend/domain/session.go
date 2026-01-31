package domain

import "time"

type Session struct {
	Id                    string     `json:"id"`
	UserId                string     `json:"user_id"`
	RefreshTokenEncrypted []byte     `json:"-"`
	UserAgent             string     `json:"user_agent"`
	IpAddress             string     `json:"ip_address"`
	DeviceName            string     `json:"device"`
	CreatedAt             time.Time  `json:"created_at"`
	LastActive            time.Time  `json:"last_active"`
	RevokedAt             *time.Time `json:"revoked_at"`
	ExpiresAt             time.Time  `json:"expires_at"`
}
