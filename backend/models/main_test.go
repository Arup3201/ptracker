package models

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/testhelpers"
	"github.com/stretchr/testify/suite"
)

type ModelTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	conn        *sql.DB
	ctx         context.Context
}

func (suite *ModelTestSuite) SetupSuite() {
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

func (suite *ModelTestSuite) TearDownSuite() {
	_, err := suite.conn.Exec("TRUNCATE users CASCADE")
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *ModelTestSuite) Cleanup(t *testing.T) {
	_, err := suite.conn.Exec("TRUNCATE projects CASCADE")
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}

func TestModel(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
