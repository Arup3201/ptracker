package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) interfaces.ProjectRepository {
	return &projectRepo{
		db: db,
	}
}

func (r *projectRepo) Create(ctx context.Context,
	name string,
	description, skills *string,
	ownerId string) (string, error) {

	id := uuid.NewString()
	project := models.Project{
		ID:          id,
		Name:        name,
		Description: description,
		Skills:      skills,
		OwnerID:     ownerId,
	}

	err := gorm.G[models.Project](r.db).Create(ctx, &project)
	if err != nil {
		return "", fmt.Errorf("gorm create: %w", err)
	}

	return project.ID, nil
}

func (r *projectRepo) Get(ctx context.Context, id string) (domain.ProjectSummary, error) {

	var err error
	var row models.ProjectSummaryRow
	err = r.db.WithContext(ctx).
		Table("projects p").
		Select(`p.id, p.name, p.description, p.skills, p.owner_id, 
				ps.unassigned_tasks, ps.ongoing_tasks, 
				ps.completed_tasks, ps.abandoned_tasks,
				p.created_at, p.updated_at
			`).
		Joins("LEFT JOIN project_summary ps ON ps.id=p.id").
		Where("p.id = ?", id).
		First(&row).
		Error
	if err == gorm.ErrRecordNotFound {
		return row.ToProjectSummaryDomain(), apierr.ErrNotFound
	} else if err != nil {
		return row.ToProjectSummaryDomain(), fmt.Errorf("gorm query: %w", err)
	}

	return row.ToProjectSummaryDomain(), nil
}

func (r *projectRepo) Delete(ctx context.Context, id string) error {

	_, err := gorm.G[domain.Project](r.db).Where(
		"id = ?",
		id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("gorm query delete: %w", err)
	}

	return nil
}
