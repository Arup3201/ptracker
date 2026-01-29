package repositories

import (
	"testing"

	"github.com/ptracker/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestProjects() {
	t := suite.T()

	t.Run("should return 2 projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p1, USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p2, USER_ONE))

		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PrivateProjects(suite.ctx, USER_ONE)

		suite.Require().NoError(err)
		expected := 2
		actual := len(projects)
		suite.Require().Equal(expected, actual)

		suite.Cleanup()
	})
}

func (suite *RepositoryTestSuite) TestPublicProjects() {
	t := suite.T()

	t.Run("should get 2 public projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p1, USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p2, USER_ONE))

		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PublicProjects(suite.ctx, USER_TWO)

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(projects))
		suite.Require().ElementsMatch(
			[]string{p1, p2},
			[]string{projects[0].Id, projects[1].Id},
		)

		suite.Cleanup()
	})
}
