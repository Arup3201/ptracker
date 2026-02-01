package repo_fixtures

import (
	"fmt"

	"github.com/ptracker/internal/domain"
)

func GetAssigneeRow(projectId, taskId, userID string) domain.Assignee {
	return domain.Assignee{
		ProjectId: projectId,
		TaskId:    taskId,
		UserId:    userID,
	}
}

func (f *Fixtures) InsertAssignee(a domain.Assignee) {
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO assignees (project_id, task_id, user_id)
		VALUES ($1,$2,$3)
		`,
		a.ProjectId,
		a.TaskId,
		a.UserId,
	)
	if err != nil {
		panic(fmt.Sprintf("insert role fixture failed: %v", err))
	}
}
