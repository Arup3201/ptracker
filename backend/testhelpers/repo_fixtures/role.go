package repo_fixtures

import (
	"fmt"

	"github.com/ptracker/domain"
)

func GetRoleRow(projectID, userID string) domain.Role {
	return domain.Role{
		ProjectId: projectID,
		UserId:    userID,
		Role:      domain.ROLE_OWNER,
	}
}

func (f *Fixtures) InsertRole(r domain.Role) {
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO roles (project_id, user_id, role)
		VALUES ($1,$2,$3)
		`,
		r.ProjectId,
		r.UserId,
		r.Role,
	)
	if err != nil {
		panic(fmt.Sprintf("insert role fixture failed: %v", err))
	}
}
