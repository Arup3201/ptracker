package repositories

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/repositories/models"
	"github.com/ptracker/internal/testhelpers/repo_fixtures"
	"gorm.io/gorm"
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
		task, _ := gorm.G[models.Task](suite.db).Where("id = ?", id).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(id, task.ID)
		suite.Require().Equal(sample_title, task.Title)
		suite.Require().Equal(sample_description, *task.Description)
		suite.Require().Equal(sample_status, task.Status.String)
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

func (suite *RepositoryTestSuite) TestGetTask() {
	t := suite.T()

	t.Run("should get title description and status of the task", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		sample_title := "sample task"
		sample_description := "sample description"
		sample_status := domain.TASK_STATUS_UNASSIGNED
		taskID := suite.fixtures.InsertTask(models.Task{
			ProjectID:   p,
			Title:       sample_title,
			Description: &sample_description,
			Status:      models.TaskStatus{String: sample_status},
		})
		repo := NewTaskRepo(suite.db)

		task, _ := repo.Get(suite.ctx, taskID)

		suite.Cleanup()

		suite.Require().Equal(taskID, task.ID)
		suite.Require().Equal(sample_title, task.Title)
		suite.Require().Equal(sample_description, *task.Description)
		suite.Require().Equal(sample_status, task.Status)
	})
	t.Run("should list assignees", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertMembership(repo_fixtures.GetMembershipRow(p, USER_TWO, domain.ROLE_MEMBER))
		taskID := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskID, USER_TWO))
		repo := NewTaskRepo(suite.db)

		_, err := repo.Get(suite.ctx, taskID)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should list 2 assignees", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertMembership(repo_fixtures.GetMembershipRow(p, USER_TWO, domain.ROLE_MEMBER))
		suite.fixtures.InsertMembership(repo_fixtures.GetMembershipRow(p, USER_THREE, domain.ROLE_MEMBER))
		taskID := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskID, USER_TWO))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskID, USER_THREE))
		repo := NewTaskRepo(suite.db)

		task, _ := repo.Get(suite.ctx, taskID)

		suite.Cleanup()

		suite.Require().Equal(2, len(task.Assignees))
	})
	t.Run("should list 2 assignees with ID", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertMembership(repo_fixtures.GetMembershipRow(p, USER_TWO, domain.ROLE_MEMBER))
		suite.fixtures.InsertMembership(repo_fixtures.GetMembershipRow(p, USER_THREE, domain.ROLE_MEMBER))
		taskID := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskID, USER_TWO))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskID, USER_THREE))
		repo := NewTaskRepo(suite.db)

		task, _ := repo.Get(suite.ctx, taskID)

		suite.Cleanup()

		suite.ElementsMatch(
			[]string{USER_TWO, USER_THREE},
			[]string{task.Assignees[0].UserID, task.Assignees[1].UserID},
		)
	})
}
