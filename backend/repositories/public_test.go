package repositories

import (
	"testing"

	"github.com/ptracker/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestPublicProjectGet() {
	t := suite.T()

	t.Run("should get public project", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_ONE))

		publicRepo := NewPublicRepo(suite.db)
		_, err := publicRepo.Get(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get public project with id", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_ONE))

		publicRepo := NewPublicRepo(suite.db)
		project, _ := publicRepo.Get(suite.ctx, p)

		suite.Cleanup()

		suite.Require().Equal(p, project.Id)
	})
}
