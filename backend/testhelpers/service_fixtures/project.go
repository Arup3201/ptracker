package service_fixtures

import "fmt"

type ProjectParams struct {
	Title       string
	Description *string
	Skills      *string
	OwnerID     string
}

func DefaultProject(ownerID string) ProjectParams {
	return ProjectParams{
		Title:   "Test Project",
		OwnerID: ownerID,
	}
}

func (f *Fixtures) Project(p ProjectParams) string {
	id, err := f.store.Project().Create(
		f.ctx,
		p.Title,
		p.Description,
		p.Skills,
		p.OwnerID,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create project fixture: %v", err))
	}

	return id
}
