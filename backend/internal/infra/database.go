package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(pgDSN string) (*gorm.DB, error) {
	var err error

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	return db, nil
}
