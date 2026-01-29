package services

import (
	"testing"

	"github.com/ptracker/domain"
	"github.com/ptracker/testhelpers/service_fixtures"
)

func (suite *ServiceTestSuite) TestJoinProject() {
	t := suite.T()

	t.Run("should join project", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewPublicService(suite.store)

		err := service.JoinProject(suite.ctx, p, USER_TWO)

		suite.Require().NoError(err)

		suite.Cleanup()
	})
	t.Run("should join project with status pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewPublicService(suite.store)

		service.JoinProject(suite.ctx, p, USER_TWO)

		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)
		suite.Require().Equal("Pending", status)

		suite.Cleanup()
	})
	t.Run("should join project with role as member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewPublicService(suite.store)

		service.JoinProject(suite.ctx, p, USER_TWO)

		var role string
		suite.db.QueryRow(
			"SELECT "+
				"role "+
				"FROM roles "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&role)
		suite.Require().Equal("Member", role)

		suite.Cleanup()
	})
	t.Run("should fail with duplicate value error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewPublicService(suite.store)
		service.JoinProject(suite.ctx, p, USER_TWO)

		err := service.JoinProject(suite.ctx, p, USER_TWO)

		suite.Require().ErrorContains(err, "duplicate value")
	})
	t.Run("should fail with transaction rollback", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		service := NewPublicService(suite.store)
		// just for testing transaction working fine...
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

		err := service.JoinProject(suite.ctx, p, USER_TWO)

		suite.Require().ErrorContains(err, "transaction: store role create")
		err = suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan()
		suite.Require().ErrorContains(err, "no rows in result set")

		suite.Cleanup()
	})
}
