package controllers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ptracker/internal/testhelpers/controller_fixtures"
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

func (suite *ControllerTestSuite) TestProjectControllerListProjects() {
	t := suite.T()

	t.Run("should list projects", func(t *testing.T) {
		client := suite.fixtures.RequestAs(USER_ONE)

		rec := client.Get("/projects")

		suite.Cleanup()

		suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	})
	t.Run("should return empty list for user with no projects", func(t *testing.T) {
		client := suite.fixtures.RequestAs(USER_TWO)

		rec := client.Get("/projects")

		suite.Cleanup()

		var response HTTPSuccessResponse[ListedPrivateProjectsResponse]
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			suite.Require().NoError(err)
		}

		suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
		suite.Require().Equal(0, len(response.Data.ProjectSummaries))
	})
}

func (suite *ControllerTestSuite) TestProjectControllerListMembers() {
	t := suite.T()

	t.Run("should get empty list of members", func(t *testing.T) {
		projectId := suite.fixtures.Project(controller_fixtures.ProjectParams{
			Name:  "Project One",
			Owner: USER_ONE,
		})
		client := suite.fixtures.RequestAs(USER_ONE)

		rec := client.Get("/projects/" + projectId + "/members")

		suite.Cleanup()

		suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)

		var response HTTPSuccessResponse[ListedMembersResponse]
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			suite.Require().NoError(err)
		}
		suite.Require().Equal(0, len(response.Data.Members))
	})
}
