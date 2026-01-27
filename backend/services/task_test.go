package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ptracker/models"
)

type store struct {
	db        *sql.DB
	userId    string
	projectId string
	taskId    string
}

func (s *store) Task() *models.TaskStore {
	return &models.TaskStore{
		DB:        s.db,
		UserId:    s.userId,
		ProjectId: s.projectId,
	}
}

func (s *store) Assignee() *models.AssigneeStore {
	return &models.AssigneeStore{
		DB:        s.db,
		ProjectId: s.projectId,
		TaskId:    s.taskId,
	}
}

func (s *store) Access() *models.AccessStore {
	return &models.AccessStore{
		DB:        s.db,
		ProjectId: s.projectId,
		UserId:    s.userId,
	}
}

func (s *store) cloneWithTx(tx *sql.Tx) *store {
	return &store{
		db: tx,
	}
}

func (s *store) WithTx(ctx context.Context, fn func(s *store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := s.cloneWithTx(tx)

	if err := fn(txStore); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (suite *ServiceTestSuite) TestPatchTask() {
	t := suite.T()

	t.Run("should patch task with user as owner", func(t *testing.T) {
		taskService := CreateTaskService(
			store,
			PROJECT_FIXTURES[0].Id,
		)
	})
	t.Run("should patch task title by owner", func(t *testing.T) {})
	t.Run("should patch task description by owner", func(t *testing.T) {})
	t.Run("should patch task with user as assignee", func(t *testing.T) {})

	t.Run("should be forbidden for patching title as assignee", func(t *testing.T) {})
	t.Run("should be forbidden for patching assignee as an assignee", func(t *testing.T) {})
	t.Run("should be forbidden for patching as just member", func(t *testing.T) {})
	t.Run("should be forbidden for patching as non member", func(t *testing.T) {})
}
