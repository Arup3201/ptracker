package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ptracker/models"
	"github.com/stretchr/testify/assert"
)

func (suite *ApiTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("success response is 200", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, res.Result().StatusCode, http.StatusOK)
	})

	t.Run("success response body is correct", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		var responseBody HTTPSuccessResponse[string]
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_SUCCESS_STATUS, responseBody.Status)
		assert.NotEqual(t, "", responseBody.Data)
	})

	t.Run("error with missing name in payload", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		var responseBody HTTPErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_ERROR_STATUS, responseBody.Status)
		assert.Equal(t, ERR_INVALID_BODY, responseBody.Error.Id)
		assert.Equal(t, "Project 'name' is missing", responseBody.Error.Message)
	})

	t.Run("error with unknown fields in payload", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak",
			"custom": "It is a custom field"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		var responseBody HTTPErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, RESPONSE_ERROR_STATUS, responseBody.Status)
		assert.Equal(t, ERR_INVALID_BODY, responseBody.Error.Id)
		assert.Equal(t, "Project must have 'name' with optional 'description' and 'skills' fields only", responseBody.Error.Message)
	})
}

func (suite *ApiTestSuite) TestProjectGet() {
	t := suite.T()

	t.Run("get project details success", func(t *testing.T) {
		payload := bytes.NewBuffer([]byte(`{
			"name": "PTracker Go", 
			"description": "Collaborative project tracking application with Go",
			"skills": "Go, React, TypeScript, PostgreSQL, Keycloak"
		}`))
		req, err := http.NewRequest("POST", "/api/projects", payload)
		if err != nil {
			log.Fatal(err)
		}
		req.AddCookie(suite.cookie)
		res := httptest.NewRecorder()
		suite.mux.ServeHTTP(res, req)
		if res.Result().StatusCode != http.StatusOK {
			t.Log("project create failed")
			t.Fail()
		}
		var createdProject HTTPSuccessResponse[string]
		if err := json.NewDecoder(res.Body).Decode(&createdProject); err != nil {
			t.Log("project create decode failed")
			t.Fail()
		}
		projectId := *(createdProject.Data)
		req, err = http.NewRequest("GET", "/api/projects/"+projectId, nil)
		if err != nil {
			t.Log("project get request create failed")
			t.Fail()
		}
		req.AddCookie(suite.cookie)
		res = httptest.NewRecorder()

		suite.mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		var projectDetails HTTPSuccessResponse[ProjectDetails]
		if err := json.NewDecoder(res.Body).Decode(&projectDetails); err != nil {
			t.Log("project get decode failed")
			t.Fail()
		}
		assert.Equal(t, RESPONSE_SUCCESS_STATUS, projectDetails.Status)
		assert.Equal(t, "PTracker Go", projectDetails.Data.Name)
		assert.Equal(t, "Collaborative project tracking application with Go", *projectDetails.Data.Description)
		assert.Equal(t, "Go, React, TypeScript, PostgreSQL, Keycloak", *projectDetails.Data.Skills)
		assert.Equal(t, models.ROLE_OWNER, projectDetails.Data.Role)
		assert.Equal(t, 0, projectDetails.Data.UnassignedTasks)
		assert.Equal(t, 0, projectDetails.Data.OngoingTasks)
		assert.Equal(t, 0, projectDetails.Data.CompletedTasks)
		assert.Equal(t, 0, projectDetails.Data.AbandonedTasks)
		assert.Equal(t, 0, projectDetails.Data.MemberCount)
	})
}
