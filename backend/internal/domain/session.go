package domain

import "time"

type Session struct {
	ID                    string
	UserID                string
	RefreshTokenEncrypted []byte
	UserAgent             string
	IpAddress             string
	DeviceName            string
	CreatedAt             time.Time
	LastActive            time.Time
	RevokedAt             *time.Time
	ExpiresAt             time.Time
}
