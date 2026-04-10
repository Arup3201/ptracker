package services

import (
	"testing"

	"github.com/ptracker/internal/repositories/models"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
	"gorm.io/gorm"
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

		joinRequest, _ := gorm.G[models.JoinRequest](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal("Pending", joinRequest.Status.String)
	})
	t.Run("should fail with duplicate value error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().Error(err, gorm.ErrDuplicatedKey)
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
}

func (suite *ServiceTestSuite) TestPublicServiceGetJoinStatus() {
	t := suite.T()

	t.Run("should get public project join status as Not Requested", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})

		status, _ := suite.publicService.GetJoinStatus(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().Equal("Not Requested", status)
	})
	t.Run("should get public project details with join status Pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		status, _ := suite.publicService.GetJoinStatus(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().Equal("Pending", status)
	})
}
