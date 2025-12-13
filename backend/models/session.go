package models

import "time"

type Session struct {
	Id               string
	UserId           string
	RefreshTokenHash string
	UserAgent        string
	IpAddress        string
	DeviceName       string
	CreatedAt        time.Time
	LastActiveAt     time.Time
	RevokedAt        time.Time
	ExpiresAt        time.Time
}
