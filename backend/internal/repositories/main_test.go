package repositories

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/internal/infra"
	"github.com/ptracker/internal/testdata"
	"github.com/ptracker/internal/testhelpers"
	"github.com/ptracker/internal/testhelpers/repo_fixtures"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *gorm.DB
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

	suite.db, err = infra.NewDatabase(suite.pgContainer.ConnectionString, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = testdata.TestMigrate(suite.db)
	if err != nil {
		log.Fatal(err)
	}

	suite.fixtures = repo_fixtures.New(suite.ctx, suite.db)

	USER_ONE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
	USER_TWO = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
	USER_THREE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	err := suite.db.Raw("TRUNCATE users CASCADE").Error
	suite.Require().NoError(err)

	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.Require().NoError(err)
	}
}

func (suite *RepositoryTestSuite) Cleanup() {
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM projects").Error
	suite.Require().NoError(err)
}

func TestRepositories(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
