package controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/testhelpers/controller_fixtures"
)

func (suite *ControllerTestSuite) TestTaskController_Update() {
	t := suite.T()

	t.Run("should update description of the task", func(t *testing.T) {
		projectId := suite.fixtures.Project(controller_fixtures.ProjectParams{
			Name:  "Project One",
			Owner: USER_ONE,
		})
		taskId := suite.fixtures.Task(controller_fixtures.TaskParams{
			ProjectID:   projectId,
			Name:        "Task one",
			Description: "",
			Status:      domain.TASK_STATUS_UNASSIGNED,
		})
		client := suite.fixtures.RequestAs(USER_ONE)
		updatedDescription := "Updated description"
		payload := UpdateTaskRequest{
			Description: &updatedDescription,
		}

		rec := client.Put(fmt.Sprintf("/projects/%s/tasks/%s", projectId, taskId), payload)

		suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)

		var description string
		suite.db.WithContext(suite.ctx).Table("tasks").Select("description").Where("id = ?", taskId).Scan(&description)

		suite.Cleanup()

		suite.Require().Equal(updatedDescription, description)
	})
}
