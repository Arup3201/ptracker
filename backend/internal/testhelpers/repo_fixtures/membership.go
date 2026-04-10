package repo_fixtures

import (
	"fmt"

	"github.com/ptracker/internal/repositories/models"
)

func GetMembershipRow(projectID, userID, role string) models.Membership {
	return models.Membership{
		ProjectID: projectID,
		UserID:    userID,
		Role:      models.UserRole{String: role},
	}
}

func (f *Fixtures) InsertMembership(m models.Membership) {
	if f.db != nil {
		if err := f.db.WithContext(f.ctx).Create(&m).Error; err != nil {
			panic(fmt.Sprintf("insert membership fixture failed: %v", err))
		}
		return
	}
}
