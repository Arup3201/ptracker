package models

import (
	"time"
)

type JoinRequest struct {
	ProjectID string `gorm:"primaryKey"`
	UserID    string `gorm:"primaryKey"`
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
