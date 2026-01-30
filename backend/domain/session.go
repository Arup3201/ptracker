package domain

import "time"

type Session struct {
	Id                    string
	UserId                string
	RefreshTokenEncrypted []byte
	UserAgent             string
	IpAddress             string
	DeviceName            string
	CreatedAt             time.Time
	LastActive            time.Time
	RevokedAt             *time.Time
	ExpiresAt             time.Time
}
