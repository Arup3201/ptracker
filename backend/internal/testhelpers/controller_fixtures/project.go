package controller_fixtures

import "github.com/ptracker/internal/domain"

type ProjectParams struct {
	Name        string
	Description *string
	Skills      *string
	Owner       string
}

func (f *ControllerFixtures) Project(project ProjectParams) string {
	id, err := f.store.Project().Create(f.ctx, project.Name, project.Description, project.Skills, project.Owner)
	if err != nil {
		panic("store project create error")
	}

	err = f.store.Role().Create(f.ctx, id, project.Owner, domain.ROLE_OWNER)
	if err != nil {
		panic("store role create error")
	}

	return id
}
