package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.NotEqual(t, "", id)
		repo.Delete(suite.ctx, id)
	})
}

func (suite *RepositoryTestSuite) TestProjectAll() {
	t := suite.T()

	t.Run("should return 2 projects", func(t *testing.T) {
		repo := NewProjectRepo(suite.db)
		projects, err := repo.All(suite.ctx, USER_ONE)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		expected := 2
		actual := len(projects)
		assert.Equal(t, expected, actual)
	})
}
