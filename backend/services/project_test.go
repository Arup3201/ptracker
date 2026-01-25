package services

import (
	"fmt"
	"testing"

	"github.com/ptracker/apierr"
	"github.com/ptracker/models"
	"github.com/stretchr/testify/assert"
)

func (suite *ServiceTestSuite) TestExploreProjects() {
	t := suite.T()
	t.Run("explore list is empty", func(t *testing.T) {
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}

		projects, err := exploreService.List(1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 0, len(projects))
	})
	t.Run("explore lists 2 projects", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			_, err := projectStore.Create(projectName, projectDescription, "C++, Python")
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		projects, err := exploreService.List(1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 2, len(projects))
		suite.Cleanup(t)
	})
	t.Run("explore lists projects where role is User", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			_, err := projectStore.Create(projectName, projectDescription, "C++, Python")
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}

		projects, err := exploreService.List(1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		for _, p := range projects {
			assert.Equal(t, "User", p.Role)
		}
		suite.Cleanup(t)
	})
	t.Run("explore lists projects where role is Owner", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			_, err := projectStore.Create(projectName, projectDescription, "C++, Python")
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}

		projects, err := exploreService.List(1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		for _, p := range projects {
			assert.Equal(t, "Owner", p.Role)
		}
		suite.Cleanup(t)
	})
}

func (suite *ServiceTestSuite) TestJoinProjectRequest() {
	projectStore := &models.ProjectStore{
		DB:     suite.conn,
		UserId: USER_FIXTURES[1].Id,
	}
	exploreService := &ProjectService{
		DB:     suite.conn,
		UserId: USER_FIXTURES[1].Id,
	}
	t := suite.T()
	t.Run("join request with pending status", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		err = exploreService.Join(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		var requestStatus string
		err = suite.conn.QueryRow("SELECT status FROM join_requests "+
			"WHERE user_id=($1) AND project_id=($2)", USER_FIXTURES[1].Id, projectId).
			Scan(&requestStatus)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, "Pending", requestStatus)
		suite.Cleanup(t)
	})
	t.Run("should return error for duplicate join request", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		projectId, err := projectStore.Create(projectName, projectDesc, projectSkills)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		err = exploreService.Join(projectId)

		err = exploreService.Join(projectId)

		assert.Equal(t, apierr.ErrDuplicate, err)
	})
}

func (suite *ServiceTestSuite) TestGetExploredProjectDetails() {
	t := suite.T()

	t.Run("should have join request status as 'Not Requested'", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}

		project, err := exploreService.Get(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, "Not Requested", project.JoinStatus)
		suite.Cleanup(t)
	})
	t.Run("should have join request status pending", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		project, err := exploreService.Get(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, "Pending", project.JoinStatus)
	})
}

func (suite *ServiceTestSuite) TestGetJoinRequests() {
	t := suite.T()

	t.Run("should have 1 join request", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		requests, err := exploreService.JoinRequests(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 1, len(requests))
		suite.Cleanup(t)
	})
	t.Run("should have join request with user", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		requests, err := exploreService.JoinRequests(projectId)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := requests[0]
		assert.Equal(t, USER_FIXTURES[1].Id, actual.User.Id)
		assert.Equal(t, USER_FIXTURES[1].Username, actual.User.Username)
		assert.Equal(t, USER_FIXTURES[1].DisplayName, actual.User.DisplayName)
		suite.Cleanup(t)
	})
}

func (suite *ServiceTestSuite) TestUpdateJoinRequestStatus() {
	t := suite.T()

	t.Run("should update join request status to accepted", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		err = exploreService.UpdateJoinRequestStatus(
			projectId,
			USER_FIXTURES[1].Id,
			"Accepted",
		)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		var status string
		err = suite.conn.QueryRow("SELECT status FROM join_requests "+
			"WHERE user_id=($1) AND project_id=($2)",
			USER_FIXTURES[1].Id, projectId).
			Scan(&status)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, status, "Accepted")
		suite.Cleanup(t)
	})
	t.Run("should update join status to rejected", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		err = exploreService.UpdateJoinRequestStatus(
			projectId,
			USER_FIXTURES[1].Id,
			"Rejected",
		)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		var status string
		err = suite.conn.QueryRow("SELECT status FROM join_requests "+
			"WHERE user_id=($1) AND project_id=($2)",
			USER_FIXTURES[1].Id, projectId).
			Scan(&status)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, status, "Rejected")
		suite.Cleanup(t)
	})
	t.Run("should give invalid value error when join status updated to pending", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		err = exploreService.UpdateJoinRequestStatus(
			projectId,
			USER_FIXTURES[1].Id,
			"Pending",
		)

		assert.Equal(t, apierr.ErrInvalidValue, err)
		suite.Cleanup(t)
	})
	t.Run("should give invalid value error for join status invalid", func(t *testing.T) {
		projectStore := &models.ProjectStore{
			DB:     suite.conn,
			UserId: USER_FIXTURES[0].Id,
		}
		projectName := fmt.Sprintf("Project %d", 1)
		projectDescription := fmt.Sprintf("Project Description %d", 1)
		projectId, err := projectStore.Create(projectName, projectDescription, "C++, Python")
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		exploreService := &ProjectService{
			DB:     suite.conn,
			UserId: USER_FIXTURES[1].Id,
		}
		err = exploreService.Join(projectId)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		err = exploreService.UpdateJoinRequestStatus(
			projectId,
			USER_FIXTURES[1].Id,
			"Invalid",
		)

		assert.Equal(t, apierr.ErrInvalidValue, err)
		suite.Cleanup(t)
	})
}
