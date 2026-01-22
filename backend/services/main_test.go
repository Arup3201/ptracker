package services

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	conn        *sql.DB
	ctx         context.Context
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnectPostgres(pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.conn = conn

	CreatFixtures(suite.conn)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	_, err := suite.conn.Exec("TRUNCATE users CASCADE")
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *ServiceTestSuite) Cleanup(t *testing.T) {
	_, err := suite.conn.Exec("TRUNCATE projects CASCADE")
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
