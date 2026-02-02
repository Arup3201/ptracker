package repositories

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestCreateTask() {
	t := suite.T()

	t.Run("should create task", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		sample_title := "sample task"
		sample_description := "sample description"
		sample_status := domain.TASK_STATUS_UNASSIGNED

		taskRepo := NewTaskRepo(suite.db)
		_, err := taskRepo.Create(suite.ctx, p,
			sample_title, sample_description, sample_status)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create task with title description and status", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		sample_title := "sample task"
		sample_description := "sample description"
		sample_status := domain.TASK_STATUS_UNASSIGNED

		taskRepo := NewTaskRepo(suite.db)
		id, _ := taskRepo.Create(suite.ctx, p,
			sample_title, sample_description, sample_status)
		var task domain.Task
		suite.db.QueryRowContext(
			suite.ctx,
			"SELECT "+
				"id, project_id, title, description, status "+
				"FROM tasks "+
				"WHERE id=($1)",
			id,
		).Scan(&task.Id, &task.ProjectId, &task.Title, &task.Description, &task.Status)

		suite.Cleanup()

		suite.Require().Equal(id, task.Id)
		suite.Require().Equal(sample_title, task.Title)
		suite.Require().Equal(sample_description, *task.Description)
		suite.Require().Equal(sample_status, task.Status)
	})
	t.Run("should create task without description", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		sample_title := "sample task"
		sample_status := domain.TASK_STATUS_UNASSIGNED

		taskRepo := NewTaskRepo(suite.db)
		_, err := taskRepo.Create(suite.ctx, p,
			sample_title, "", sample_status)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
}
