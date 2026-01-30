package repo_fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/domain"
)

func RandomTaskRow(projectId, status string) domain.Task {
	tId := uuid.NewString()

	desc := "Description " + tId

	return domain.Task{
		Id:          tId,
		ProjectId:   projectId,
		Title:       "Test Task " + tId,
		Description: &desc,
		Status:      status,
	}
}

func (f *Fixtures) InsertTask(t domain.Task) string {
	now := time.Now()
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO tasks (
			id, project_id, title, description, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$6)
		`,
		t.Id,
		t.ProjectId,
		t.Title,
		t.Description,
		t.Status,
		now,
	)
	if err != nil {
		panic(fmt.Sprintf("insert task fixture failed: %v", err))
	}

	return t.Id
}
