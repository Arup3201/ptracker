package repo_fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/internal/domain"
)

func RandomProjectRow(ownerID string) domain.Project {
	pId := uuid.NewString()

	desc := "Description " + pId
	skills := "C++, Python"

	return domain.Project{
		Id:          pId,
		Name:        "Test Project" + pId,
		Description: &desc,
		Skills:      &skills,
		Owner:       ownerID,
	}
}

func (f *Fixtures) InsertProject(p domain.Project) string {
	now := time.Now()
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO projects (
			id, name, description, skills, owner, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$6)
		`,
		p.Id,
		p.Name,
		p.Description,
		p.Skills,
		p.Owner,
		now,
	)
	if err != nil {
		panic(fmt.Sprintf("insert project fixture failed: %v", err))
	}

	f.InsertRole(GetRoleRow(p.Id, p.Owner, domain.ROLE_OWNER))

	return p.Id
}
