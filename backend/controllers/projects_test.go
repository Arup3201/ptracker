package controllers

import (
	"encoding/json"
	"net/http"
	"testing"
)

func (suite *ControllerTestSuite) TestProjectControllerCreate() {
	t := suite.T()

	t.Run("should create project", func(t *testing.T) {
		test_name := "Test Project"
		test_description := "Test Description"
		project := CreateProjectRequest{
			Name:        test_name,
			Description: &test_description,
		}
		client := suite.fixtures.RequestAs(USER_ONE)

		rec := client.Post("/projects", project)

		suite.Cleanup()

		suite.Require().Equal(http.StatusCreated, rec.Result().StatusCode)
	})
	t.Run("should create project with name description and owner", func(t *testing.T) {
		test_name := "Test Project"
		test_description := "Test Description"
		project := CreateProjectRequest{
			Name:        test_name,
			Description: &test_description,
		}
		client := suite.fixtures.RequestAs(USER_ONE)

		rec := client.Post("/projects", project)

		var response HTTPSuccessResponse[string]
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			suite.Require().NoError(err)
		}

		var name, description, owner string
		err := suite.db.QueryRow(
			"SELECT name, description, owner FROM projects WHERE id=($1)",
			response.Data,
		).Scan(&name, &description, &owner)
		suite.Require().NoError(err)

		suite.Cleanup()

		suite.Require().Equal(test_name, name)
		suite.Require().Equal(test_description, description)
		suite.Require().Equal(USER_ONE, owner)
	})
	t.Run("should get bad request without name", func(t *testing.T) {
		test_description := "Test Description"
		project := CreateProjectRequest{
			Description: &test_description,
		}
		client := suite.fixtures.RequestAs(USER_ONE)

		rec := client.Post("/projects", project)

		suite.Cleanup()

		var response HTTPErrorResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			suite.Require().NoError(err)
		}

		suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
		suite.Require().Equal(response.Error.Message, "Project 'name' is missing")
	})
	t.Run("should get internal error invalid user", func(t *testing.T) {
		test_name := "Test Project"
		test_description := "Test Description"
		project := CreateProjectRequest{
			Name:        test_name,
			Description: &test_description,
		}
		client := suite.fixtures.RequestAs("abcfd")

		rec := client.Post("/projects", project)

		suite.Cleanup()

		suite.Require().Equal(http.StatusInternalServerError, rec.Result().StatusCode)
	})
}
