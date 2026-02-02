package service_fixtures

import (
	"fmt"
)

type TaskParams struct {
	ProjectId   string
	Title       string
	Description string
	Status      string
	Assignees   []string
}

func (f *Fixtures) Task(t TaskParams) string {
	id, err := f.store.Task().Create(
		f.ctx,
		t.ProjectId,
		t.Title,
		t.Description,
		t.Status,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create task fixture: %v", err))
	}

	for _, assignee := range t.Assignees {
		err = f.store.Assignee().Create(f.ctx, t.ProjectId, id, assignee)
		if err != nil {
			panic(fmt.Sprintf("failed to add assignee fixture: %v", err))
		}
	}

	return id
}
