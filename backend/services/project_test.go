package services

import (
	"testing"
)

func (suite *ServiceTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("should create project with user one as owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Require().NoError(err)
		suite.Require().NotEqual("", id)

		service.store.Project().Delete(suite.ctx, id)
	})
}

func (suite *ServiceTestSuite) TestGetPrivateProject() {
	t := suite.T()

	t.Run("should return project correctly", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := service.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		suite.Require().Equal(id, project.Id)

		service.store.Project().Delete(suite.ctx, id)
	})
}
