package repo_fixtures

import (
	"fmt"

	"github.com/ptracker/domain"
)

func GetJoinRequest(projectID, userID string) domain.JoinRequest {
	return domain.JoinRequest{
		ProjectId: projectID,
		UserId:    userID,
		Status:    domain.JOIN_STATUS_PENDING,
	}
}

func (f *Fixtures) InsertJoinRequest(j domain.JoinRequest) {
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO join_requests (project_id, user_id, status)
		VALUES ($1,$2,$3)
		`,
		j.ProjectId,
		j.UserId,
		j.Status,
	)
	if err != nil {
		panic(fmt.Sprintf("insert join request fixture failed: %v", err))
	}
}
