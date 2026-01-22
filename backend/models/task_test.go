package models

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *ModelTestSuite) TestTaskCreate() {
	projectStore := &ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[0].Id,
	}

	var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
	projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
	if err != nil {
		log.Fatal(err)
	}

	taskStore := &TaskStore{
		DB:        suite.conn,
		UserId:    USER_FIXTURES[0].Id,
		ProjectId: projectId,
	}

	t := suite.T()
	t.Run("get project tasks success", func(t *testing.T) {
		taskTitle, taskDescription, taskStatus := "Test Task", "Test Description", "Ongoing"

		taskId, err := taskStore.Create(taskTitle, taskDescription, taskStatus)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.NotEqual(t, "", taskId)
		suite.Cleanup(t)
	})
}

func (suite *ModelTestSuite) TestTaskAll() {
	projectStore := &ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[0].Id,
	}
	t := suite.T()
	t.Run("get project tasks success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		taskStore := &TaskStore{
			DB:        suite.conn,
			UserId:    USER_FIXTURES[0].Id,
			ProjectId: projectId,
		}
		for i := range 10 {
			_, err := taskStore.Create(fmt.Sprintf("Test task title %d", i), fmt.Sprintf("Test task description %d", i), "Ongoing")
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		tasks, err := taskStore.All()

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 10, len(tasks))
		suite.Cleanup(t)
	})
}
