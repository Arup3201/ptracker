package repositories

import (
	"testing"

	"github.com/ptracker/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestProjectCreate() {
	t := suite.T()

	t.Run("should create project correctly", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		repo := NewProjectRepo(suite.db)
		id, err := repo.Create(suite.ctx,
			sample_name, &sample_description, &sample_skills,
			USER_ONE)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		suite.Require().NotEqual(t, "", id)

		repo.Delete(suite.ctx, id)
	})
}

func (suite *RepositoryTestSuite) TestProjectAll() {
	t := suite.T()

	t.Run("should return 2 projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p1, USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p2, USER_ONE))
		repo := NewProjectRepo(suite.db)
		projects, err := repo.All(suite.ctx, USER_ONE)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		expected := 2
		actual := len(projects)
		suite.Require().Equal(expected, actual)

		repo.Delete(suite.ctx, p1)
		repo.Delete(suite.ctx, p2)
	})
}
