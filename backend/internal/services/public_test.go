package services

import (
	"testing"

	"github.com/ptracker/internal/testhelpers/service_fixtures"
)

func (suite *ServiceTestSuite) TestJoinProject() {
	t := suite.T()

	t.Run("should join project", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})

		err := suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should join project with status pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})

		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		var status string
		suite.db.QueryRowContext(
			suite.ctx,
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal("Pending", status)
	})
	t.Run("should fail with duplicate value error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "duplicate value")
	})
}

func (suite *ServiceTestSuite) TestPublicServiceGet() {
	t := suite.T()

	t.Run("should get public project details", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})

		_, err := suite.publicService.GetPublicProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get public project details with join status Not Requested", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})

		project, _ := suite.publicService.GetPublicProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().Equal(project.JoinStatus, "Not Requested")
	})
	t.Run("should get public project details with join status Pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		project, _ := suite.publicService.GetPublicProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().Equal(project.JoinStatus, "Pending")
	})
}
