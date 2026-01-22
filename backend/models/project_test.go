package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *ModelTestSuite) TestProjectStoreCreate() {
	projectStore := &ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[0].Id,
	}
	t := suite.T()
	t.Run("create project success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"

		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.NotEqual(t, projectId, "")
		var p Project
		suite.conn.QueryRow(
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			projectId,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		assert.Equal(t, "Test Project 1", p.Name)
		assert.Equal(t, "Test Project Description", *p.Description)
		assert.Equal(t, "Python", *p.Skills)
		assert.Equal(t, USER_FIXTURES[0].Id, p.Owner)
		suite.Cleanup(t)
	})
}

func (suite *ModelTestSuite) TestProjectStoreGet() {
	projectStore := &ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[0].Id,
	}
	t := suite.T()
	t.Run("get project details success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		details, err := projectStore.Get(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := *details
		assert.Equal(t, projectName, actual.Name)
		assert.Equal(t, projectDesc, *actual.Description)
		assert.Equal(t, projectSkills, *actual.Skills)
		assert.Equal(t, USER_FIXTURES[0].Id, actual.Owner)
		assert.Equal(t, 0, actual.UnassignedTasks)
		assert.Equal(t, 0, actual.OngoingTasks)
		assert.Equal(t, 0, actual.CompletedTasks)
		assert.Equal(t, 0, actual.AbandonedTasks)
		suite.Cleanup(t)
	})
}
