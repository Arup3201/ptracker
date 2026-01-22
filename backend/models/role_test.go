package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *ModelTestSuite) TestCanAccess() {
	projectStore := &ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[0].Id,
	}
	t := suite.T()
	t.Run("can access", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		roleStore := &RoleStore{
			DB:        suite.conn,
			UserId:    USER_FIXTURES[0].Id,
			ProjectId: projectId,
		}
		access, err := roleStore.CanAccess()

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, true, access)
		suite.Cleanup(t)
	})

	t.Run("can't access", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		roleStore := &RoleStore{
			DB:        suite.conn,
			UserId:    USER_FIXTURES[1].Id,
			ProjectId: projectId,
		}

		access, err := roleStore.CanAccess()

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, false, access)
		suite.Cleanup(t)
	})
}
