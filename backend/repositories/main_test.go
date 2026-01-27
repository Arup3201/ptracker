package repositories

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *sql.DB
	ctx         context.Context
}

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

	CreatFixtures(suite.db)
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
