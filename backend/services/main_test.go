package services

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/db"
	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
	"github.com/ptracker/testhelpers"
	"github.com/ptracker/testhelpers/fixtures"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *testhelpers.PostgresContainer
	db          *sql.DB
	store       stores.Store
	fixtures    *fixtures.Fixtures
}

var USER_ONE, USER_TWO string
var PROJECT_ONE, PROJECT_TWO string

func (suite *ServiceTestSuite) SetupSuite() {
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

	suite.store = NewSQLStore(suite.db)
	suite.fixtures = fixtures.New(suite.ctx, suite.store)

	USER_ONE = suite.fixtures.User(fixtures.UserParams{
		IDPSubject:  "sub-234",
		IDPProvider: "facebook",
		Username:    "bob",
		Email:       "bob@example.com",
	})
	USER_TWO = suite.fixtures.User(fixtures.UserParams{
		IDPSubject:  "sub-345",
		IDPProvider: "twitter",
		Username:    "alice",
		Email:       "alice@example.com",
	})
	PROJECT_ONE = suite.fixtures.Project(fixtures.ProjectParams{
		Title:   "Project Fixture A",
		OwnerID: USER_ONE,
	})
	suite.fixtures.Role(PROJECT_ONE, USER_ONE, domain.ROLE_OWNER)
	PROJECT_TWO = suite.fixtures.Project(fixtures.ProjectParams{
		Title:   "Project Fixture B",
		OwnerID: USER_TWO,
	})
	suite.fixtures.Role(PROJECT_TWO, USER_TWO, domain.ROLE_OWNER)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	_, err := suite.db.Exec("TRUNCATE users CASCADE")
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
