package repositories

import (
	"testing"

	"github.com/ptracker/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestProjectAll() {
	t := suite.T()

	t.Run("should return 2 projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p1, USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p2, USER_ONE))
		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PrivateProjects(suite.ctx, USER_ONE)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		expected := 2
		actual := len(projects)
		suite.Require().Equal(expected, actual)

		projectRepo := NewProjectRepo(suite.db)
		projectRepo.Delete(suite.ctx, p1)
		projectRepo.Delete(suite.ctx, p2)
	})
}
