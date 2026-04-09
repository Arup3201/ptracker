package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type MembershipRepo struct {
	db *gorm.DB
}

func NewMembershipRepo(db *gorm.DB) interfaces.MembershipRepository {
	return &MembershipRepo{
		db: db,
	}
}

func (r *MembershipRepo) Create(ctx context.Context,
	projectId, userId, userRole string) error {

	role := models.Membership{
		ProjectID: projectId,
		UserID:    userId,
		Role:      userRole,
	}
	err := gorm.G[models.Membership](r.db).Create(ctx, &role)
	if err != nil {
		return fmt.Errorf("gorm create: %w", err)
	}

	return nil
}

func (r *MembershipRepo) Role(ctx context.Context, projectId, userId string) (string, error) {

	membership, err := gorm.G[models.Membership](r.db).Where(
		"project_id = ? AND user_id = ?",
		projectId, userId).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return "", apierr.ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("gorm query: %w", err)
	}

	return membership.Role, nil
}

func (r *MembershipRepo) CountMembers(ctx context.Context, projectId string) (int64, error) {

	var cnt int64
	err := r.db.Model(&models.Membership{}).Where("project_id = ?", projectId).Count(&cnt).Error
	if err != nil {
		return -1, fmt.Errorf("gorm db count query: %w", err)
	}

	return cnt, nil
}
