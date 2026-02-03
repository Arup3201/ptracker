package repo_fixtures

import (
	"fmt"
	"time"

	"github.com/ptracker/internal/domain"
)

func GetRoleRow(projectID, userID, role string) domain.Role {
	return domain.Role{
		ProjectId: projectID,
		UserId:    userID,
		Role:      role,
	}
}

func (f *Fixtures) InsertRole(r domain.Role) {
	now := time.Now()
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO roles (project_id, user_id, role, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$4)
		`,
		r.ProjectId,
		r.UserId,
		r.Role,
		now,
	)
	if err != nil {
		panic(fmt.Sprintf("insert role fixture failed: %v", err))
	}
}
