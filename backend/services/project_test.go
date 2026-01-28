package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *ServiceTestSuite) TestCreateProject() {
	t := suite.T()

	fakeStore := &fakeStore{}

	t.Run("should create project with user one as owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(fakeStore)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_FIXTURES[0],
		)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.NotEqual(t, "", id)
	})
}
