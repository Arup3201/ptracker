package repositories

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/testhelpers"
	"github.com/ptracker/testhelpers/repo_fixtures"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *sql.DB
	fixtures    *repo_fixtures.Fixtures
	ctx         context.Context
}

var USER_ONE, USER_TWO string
var PROJECT_ONE, PROJECT_TWO, PROJECT_THREE string

func (suite *RepositoryTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	suite.pgContainer, err = testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.db, err = db.ConnectPostgres(suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.fixtures = repo_fixtures.New(suite.ctx, suite.db)

	USER_ONE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
	USER_TWO = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())

	PROJECT_ONE = suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
	suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(PROJECT_ONE, USER_ONE))

	PROJECT_TWO = suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
	suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(PROJECT_TWO, USER_ONE))

	PROJECT_THREE = suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_TWO))
	suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(PROJECT_THREE, USER_TWO))
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	_, err := suite.db.Exec("TRUNCATE users CASCADE")
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func TestModel(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
