package models

import (
	"time"
)

type Membership struct {
	ProjectID string `gorm:"primaryKey"`
	UserID    string `gorm:"primaryKey"`
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
