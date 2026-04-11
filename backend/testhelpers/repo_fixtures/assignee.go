package repo_fixtures

import (
	"fmt"

	"github.com/ptracker/core/models"
)

func GetAssigneeRow(projectID, taskID, userID string) models.Assignee {
	return models.Assignee{
		ProjectID: projectID,
		TaskID:    taskID,
		UserID:    userID,
	}
}

func (f *Fixtures) InsertAssignee(a models.Assignee) {
	if f.db != nil {
		if err := f.db.WithContext(f.ctx).Create(&a).Error; err != nil {
			panic(fmt.Sprintf("insert assignee fixture failed: %v", err))
		}
		return
	}
}
