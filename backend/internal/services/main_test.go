package services

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/ptracker/internal/db"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/testhelpers"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *testhelpers.PostgresContainer
	db          *sql.DB
	store       interfaces.Store
	fixtures    *service_fixtures.Fixtures
}

var USER_ONE, USER_TWO, USER_THREE string

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

	redisContainer, err := testhelpers.CreateRedisContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := redisContainer.ConnectionString(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := redis.ParseURL(connString)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(opt)
	redis := db.NewRedisInMemory(redisClient)
	rateLimiter := db.NewRedisRateLimiter(redisClient, 5, 3)
	suite.store = NewStorage(suite.db, redis, rateLimiter)
	suite.fixtures = service_fixtures.New(suite.ctx, suite.store)

	USER_ONE = suite.fixtures.User(service_fixtures.UserParams{
		IDPSubject:  "sub-234",
		IDPProvider: "facebook",
		Username:    "bob",
		Email:       "bob@example.com",
	})
	USER_TWO = suite.fixtures.User(service_fixtures.UserParams{
		IDPSubject:  "sub-345",
		IDPProvider: "twitter",
		Username:    "alice",
		Email:       "alice@example.com",
	})
	USER_THREE = suite.fixtures.User(service_fixtures.UserParams{
		IDPSubject:  "sub-456",
		IDPProvider: "twitter",
		Username:    "mevis",
		Email:       "mevis@example.com",
	})
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

func (suite *ServiceTestSuite) Cleanup() {
	_, err := suite.db.Exec("DELETE FROM projects")
	suite.Require().NoError(err)
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
