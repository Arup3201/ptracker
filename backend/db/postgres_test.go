package db

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	IDPProvider     = "keycloak"
	IDPSubject      = "f6e1d9a0-7b3c-4d5e-8f2a-1c9b8e7d6f5a"
	TestKCRealm     = "ptracker"
	TestUsername    = "test_user"
	TestDisplayName = "Test User"
	TestEmail       = "test@example.com"
)

type PGTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
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
}

func (suite *PGTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *PGTestSuite) TestCreateProject() {
	user, err := CreateUser(IDPSubject, IDPProvider, TestUsername, TestDisplayName, TestEmail, "")
	if err != nil {
		log.Fatal(err)
	}

	t := suite.T()
	t.Run("create project success", func(t *testing.T) {
		var projectName, projectDesc, projectSkills = "Test Project 1", "Test Project Description", "Python"

		project, err := CreateProject(projectName, projectDesc, projectSkills, user.Id)

		if err != nil {
			t.Fail()
			t.Log(err)
		}
		actual := *project
		assert.Equal(t, projectName, actual.Name)
		assert.Equal(t, projectDesc, *actual.Description)
		assert.Equal(t, projectSkills, *actual.Skills)
		assert.Equal(t, user.Id, actual.Owner)
	})
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
