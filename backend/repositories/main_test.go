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

var USER_ONE, USER_TWO, USER_THREE string

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
	USER_THREE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	_, err := suite.db.Exec("TRUNCATE users CASCADE")
	suite.Require().NoError(err)

	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.Require().NoError(err)
	}
}

func (suite *RepositoryTestSuite) Cleanup() {
	_, err := suite.db.Exec("DELETE FROM projects")
	suite.Require().NoError(err)
}

func TestRepositories(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
