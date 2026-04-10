package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type PublicRepo struct {
	db *gorm.DB
}

func NewPublicRepo(db *gorm.DB) interfaces.PublicRepository {
	return &PublicRepo{
		db: db,
	}
}

func (r *PublicRepo) Get(ctx context.Context, projectId string) (domain.ProjectPublicDetail, error) {

	var row models.ProjectPublicDetailRow
	err := r.db.WithContext(ctx).
		Table("projects p").
		Select(`p.id, p.name, p.description, p.skills, 
			u.id as owner_id, 
			u.username as owner_username, 
			u.display_name  as owner_display_name, 
			u.email as owner_email, 
			u.avatar_url as owner_avatar_url, 
			ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, 
			p.created_at, p.updated_at`).
		Joins("LEFT JOIN project_summary as ps ON p.id=ps.id").
		Joins("INNER JOIN users as u ON u.id=p.owner_id").
		Where("p.id = ?", projectId).
		Scan(&row).Error
	if err != nil {
		return row.ToProjectPublicDetailDomain(), fmt.Errorf("db query row context: %w", err)
	}

	return row.ToProjectPublicDetailDomain(), nil
}
