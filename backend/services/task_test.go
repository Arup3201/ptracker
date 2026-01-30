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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
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

		suite.Require().Equal("Unassigned", status)
	})
	t.Run("should be forbidden for user two", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewTaskService(suite.store)
		sample_title := ""
		sample_description := "sample description"

		_, err := service.CreateTask(suite.ctx, p, sample_title, &sample_description, USER_ONE)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "invalid value")
	})
}
