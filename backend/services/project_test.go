package services

import (
	"testing"

	"github.com/ptracker/domain"
	"github.com/ptracker/testhelpers/fake"
	"github.com/stretchr/testify/assert"
)

func (suite *ServiceTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("should create project with user one as owner", func(t *testing.T) {
		store := fake.NewStore()
		store.User().(*fake.UserRepo).
			WithUser(domain.User{Id: "u1", Username: "Alice"})
		store.Project().(*fake.ProjectRepo).
			WithProject(domain.Project{Id: "p1", Name: "Project A"})
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			"u1",
		)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.NotEqual(t, "", id)
	})
}

func (suite *ServiceTestSuite) TestGetPrivateProject() {
	t := suite.T()

	t.Run("should return project correctly", func(t *testing.T) {
		store := fake.NewStore()
		store.User().(*fake.UserRepo).
			WithUser(domain.User{Id: "u1", Username: "Alice"})
		store.Project().(*fake.ProjectRepo).
			WithProject(domain.Project{Id: "p1", Name: "Project A"})
		service := NewProjectService(store)

		project, err := service.GetPrivateProject(suite.ctx, "p1", "u1")

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, "p1", project.Id)
		assert.Equal(t, "Project A", project.Name)
	})
}
