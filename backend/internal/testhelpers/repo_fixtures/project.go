package repo_fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/repositories/models"
)

func RandomProjectRow(ownerID string) models.Project {
	pId := uuid.NewString()

	desc := "Description " + pId
	skills := "C++, Python"

	return models.Project{
		ID:          pId,
		Name:        "Test Project" + pId,
		Description: &desc,
		Skills:      &skills,
		OwnerID:     ownerID,
	}
}

func (f *Fixtures) InsertProject(p models.Project) string {
	if f.db != nil {
		if err := f.db.WithContext(f.ctx).Create(&p).Error; err != nil {
			panic(fmt.Sprintf("insert project fixture failed: %v", err))
		}
		f.InsertMembership(GetMembershipRow(p.ID, p.OwnerID, domain.ROLE_OWNER))
		return p.ID
	}
	return ""
}
