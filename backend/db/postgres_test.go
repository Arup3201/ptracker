package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ptracker/models"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	IDPProvider     = "keycloak"
	IDPSubject      = "f6e1d9a0-7b3c-4d5e-8f2a-1c9b8e7d6f5a"
	IDPSubject_1    = "f6e1d9a0-7b3c-4d5e-8f2a-1c9b8e7d6f5b"
	TestKCRealm     = "ptracker"
	TestUsername    = "test_user"
	TestDisplayName = "Test User"
	TestEmail       = "test@example.com"
)

type PGTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	user        models.User
	ctx         context.Context
}

func (suite *PGTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = ConnectPostgres(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer

	user, err := CreateUser(IDPSubject, IDPProvider, TestUsername, TestDisplayName, TestEmail, "")
	if err != nil {
		log.Fatal(err)
	}

	suite.user = *user
}

func (suite *PGTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *PGTestSuite) Cleanup(t *testing.T) {
	err := truncateTable("projects")
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}

func (suite *PGTestSuite) TestCreateProject() {

	t := suite.T()
	t.Run("create project success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"

		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := *project
		assert.Equal(t, projectName, actual.Name)
		assert.Equal(t, projectDesc, *actual.Description)
		assert.Equal(t, projectSkills, *actual.Skills)
		assert.Equal(t, suite.user.Id, actual.Owner.Id)
		suite.Cleanup(t)
	})
}

func (suite *PGTestSuite) TestCanAccess() {
	t := suite.T()
	t.Run("can access", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		access, err := CanAccess(suite.user.Id, project.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, true, access)
		suite.Cleanup(t)
	})

	t.Run("can't access", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		user, err := CreateUser("some random subject", IDPProvider, "Test 2", "Test 2", "test2@example.com", "")
		if err != nil {
			log.Fatal(err)
		}

		access, err := CanAccess(user.Id, project.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, false, access)
		suite.Cleanup(t)
	})
}

func (suite *PGTestSuite) TestGetProjectDetails() {
	t := suite.T()
	t.Run("get project details success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)
		if err != nil {
			t.Fail()
			t.Log(err)
		}

		details, err := GetProjectDetails(suite.user.Id, project.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := *details
		assert.Equal(t, projectName, actual.Name)
		assert.Equal(t, projectDesc, *actual.Description)
		assert.Equal(t, projectSkills, *actual.Skills)
		assert.Equal(t, suite.user.Id, actual.Owner.Id)
		assert.Equal(t, suite.user.Username, actual.Owner.Username)
		assert.Equal(t, suite.user.DisplayName, actual.Owner.DisplayName)
		assert.Equal(t, models.ROLE_OWNER, actual.Role)
		assert.Equal(t, 0, actual.UnassignedTasks)
		assert.Equal(t, 0, actual.OngoingTasks)
		assert.Equal(t, 0, actual.CompletedTasks)
		assert.Equal(t, 0, actual.AbandonedTasks)
		assert.Equal(t, 0, actual.MemberCount)
		suite.Cleanup(t)
	})
}

func (suite *PGTestSuite) GetCreateTask() {
	t := suite.T()
	t.Run("get project tasks success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		taskTitle, taskDescription, taskStatus := "Test Task", "Test Description", "Ongoing"

		task, err := CreateProjectTask(taskTitle, taskDescription, taskStatus, project.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := *task
		assert.Equal(t, taskTitle, actual.Title)
		assert.Equal(t, taskDescription, *actual.Description)
		assert.Equal(t, taskStatus, actual.Status)
		suite.Cleanup(t)
	})
}

func (suite *PGTestSuite) TestGetProjectTasks() {
	t := suite.T()
	t.Run("get project tasks success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"
		project, err := CreateProject(projectName, projectDesc, projectSkills, suite.user.Id)
		if err != nil {
			t.Fail()
			t.Log(err)
		}
		for i := range 10 {
			_, err := CreateProjectTask(fmt.Sprintf("Test task title %d", i), fmt.Sprintf("Test task description %d", i), "Ongoing", project.Id)
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		tasks, err := GetProjectTasks(project.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 10, len(tasks))
		suite.Cleanup(t)
	})
}

func (suite *PGTestSuite) TestExploreProjects() {
	t := suite.T()

	user, err := CreateUser(IDPSubject_1, IDPProvider, "user_2", "Test User 2", "test2@example.com", "")
	if err != nil {
		log.Fatal(err)
	}

	t.Run("explore list is empty", func(t *testing.T) {
		projects, err := GetExploredProjects(user.Id, 1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 0, len(projects))
	})
	t.Run("explore lists 2 projects", func(t *testing.T) {
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			_, err := CreateProject(projectName, projectDescription, "C++, Python", user.Id)
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		projects, err := GetExploredProjects(user.Id, 1, 10)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		assert.Equal(t, 2, len(projects))
		suite.Cleanup(t)
	})
	t.Run("explore lists projects where role is User", func(t *testing.T) {
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			CreateProject(projectName, projectDescription, "C++, Python", suite.user.Id)
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		projects, err := GetExploredProjects(user.Id, 1, 10)

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
		for i := range 2 {
			projectName := fmt.Sprintf("Project %d", i+1)
			projectDescription := fmt.Sprintf("Project Description %d", i+1)
			CreateProject(projectName, projectDescription, "C++, Python", user.Id)
			if err != nil {
				t.Fail()
				t.Log(err)
			}
		}

		projects, err := GetExploredProjects(user.Id, 1, 10)

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

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
