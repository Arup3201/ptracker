package services

import (
	"testing"

	"github.com/ptracker/domain"
	"github.com/ptracker/testhelpers/service_fixtures"
)

func (suite *ServiceTestSuite) TestCreateTask() {
	t := suite.T()

	t.Run("should create task", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		service := NewTaskService(suite.store)
		sample_title := "sample task"
		sample_description := "sample description"

		_, err := service.CreateTask(suite.ctx, p, sample_title, &sample_description, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create task with unassigned status", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		service := NewTaskService(suite.store)
		sample_title := "sample task"
		sample_description := "sample description"

		taskId, _ := service.CreateTask(suite.ctx, p, sample_title, &sample_description, USER_ONE)
		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM tasks "+
				"WHERE id=($1)",
			taskId,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(domain.TASK_STATUS_UNASSIGNED, status)
	})
	t.Run("should be forbidden for user two", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)
		service := NewTaskService(suite.store)
		sample_title := "sample task"
		sample_description := "sample description"

		_, err := service.CreateTask(suite.ctx, p, sample_title, &sample_description, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be invalid with empty task title", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		service := NewTaskService(suite.store)
		sample_title := ""
		sample_description := "sample description"

		_, err := service.CreateTask(suite.ctx, p, sample_title, &sample_description, USER_ONE)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "invalid value")
	})
}

func (suite *ServiceTestSuite) TestListTasks() {
	t := suite.T()

	t.Run("should get list of tasks", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		service := NewTaskService(suite.store)

		_, err := service.ListTasks(suite.ctx, p, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get list of tasks for member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)
		service := NewTaskService(suite.store)

		_, err := service.ListTasks(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get forbidden error for non-member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		service := NewTaskService(suite.store)

		_, err := service.ListTasks(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}

func (suite *ServiceTestSuite) TestGetTask() {
	t := suite.T()

	t.Run("should get task", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		service := NewTaskService(suite.store)

		_, err := service.GetTask(suite.ctx, projectId, taskId, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get task with id title status", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     sample_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		service := NewTaskService(suite.store)

		task, _ := service.GetTask(suite.ctx, projectId, taskId, USER_ONE)

		suite.Cleanup()

		suite.Require().Equal(taskId, task.Id)
		suite.Require().Equal(sample_title, task.Title)
		suite.Require().Equal(domain.TASK_STATUS_UNASSIGNED, task.Status)
	})
	t.Run("should be forbidden for non-member", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     sample_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		service := NewTaskService(suite.store)

		_, err := service.GetTask(suite.ctx, projectId, taskId, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}
