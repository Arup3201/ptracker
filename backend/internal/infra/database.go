package infra

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(pgDSN string, config *gorm.Config) (*gorm.DB, error) {
	var err error

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(pgDSN), config)
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	return db, nil
}
