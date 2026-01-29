package repositories

import (
	"testing"
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
