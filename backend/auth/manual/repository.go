package manual

import (
	"context"
	"fmt"

	"github.com/ptracker/models"
	"gorm.io/gorm"
)

type ManualAccountRepository struct {
	db *gorm.DB
}

func NewManualAccountRepository(db *gorm.DB) *ManualAccountRepository {
	return &ManualAccountRepository{
		db: db,
	}
}

func (r *ManualAccountRepository) WithTx(tx *gorm.DB) *ManualAccountRepository {
	return NewManualAccountRepository(tx)
}

func (r *ManualAccountRepository) Create(ctx context.Context,
	userID, email string,
	passwordHash []byte,
	emailVerified bool) error {

	manualAccount := models.ManualAccount{
		UserID:        userID,
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: emailVerified,
	}
	err := gorm.G[models.ManualAccount](r.db).Create(ctx, &manualAccount)
	if err != nil {
		return fmt.Errorf("gorm create: %w", err)
	}

	return nil
}
