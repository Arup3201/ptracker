package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type JoinRequestRepo struct {
	db *gorm.DB
}

func NewJoinRequestRepo(db *gorm.DB) interfaces.JoinRequestRepository {
	return &JoinRequestRepo{
		db: db,
	}
}

func (r *JoinRequestRepo) Create(ctx context.Context,
	projectId, userId, joinStatus string) error {

	var err error

	joinReq := models.JoinRequest{
		ProjectID: projectId,
		UserID:    userId,
		Status: models.JoinStatus{
			String: joinStatus,
		},
	}
	err = gorm.G[models.JoinRequest](r.db).Create(ctx, &joinReq)
	if err == gorm.ErrDuplicatedKey {
		return apierr.ErrDuplicate
	} else if err != nil {
		return fmt.Errorf("gorm create: %w", err)
	}

	return nil
}

func (r *JoinRequestRepo) Status(ctx context.Context, projectId, userId string) (string, error) {

	joinReq, err := gorm.G[models.JoinRequest](r.db).Where(
		"project_id = ? AND user_id = ?",
		projectId, userId).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return "", apierr.ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("gorm query: %w", err)
	}

	return joinReq.Status.String, nil
}

func (r *JoinRequestRepo) Update(ctx context.Context, projectId, userId, joinStatus string) error {

	joinReq, err := gorm.G[models.JoinRequest](r.db).Where(
		"project_id = ? AND user_id = ?",
		projectId, userId).First(ctx)
	if err != nil {
		return fmt.Errorf("join request get: %w", err)
	}

	joinReq.Status = models.JoinStatus{
		String: joinStatus,
	}
	err = r.db.Save(&joinReq).Error
	if err == gorm.ErrInvalidValue {
		return apierr.ErrInvalidValue
	} else if err != nil {
		return fmt.Errorf("gorm db save: %w", err)
	}

	return nil
}
