package models

import (
	"time"

	"github.com/ptracker/internal/domain"
)

type Session struct {
	ID                    string `gorm:"primaryKey"`
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

func (s Session) ToSessionDomain() domain.Session {
	return domain.Session{
		ID:                    s.ID,
		UserID:                s.UserID,
		RefreshTokenEncrypted: s.RefreshTokenEncrypted,
		UserAgent:             s.UserAgent,
		IpAddress:             s.IpAddress,
		DeviceName:            s.DeviceName,
		CreatedAt:             s.CreatedAt,
		LastActive:            s.LastActive,
		RevokedAt:             s.RevokedAt,
		ExpiresAt:             s.ExpiresAt,
	}
}
