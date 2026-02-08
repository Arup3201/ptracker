package services

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/infra"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/testhelpers"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type mockNotifier struct{}

func (n *mockNotifier) Notify(ctx context.Context,
	user string, message domain.Message) error {
	return nil
}

func (n *mockNotifier) BatchNotify(ctx context.Context,
	users []string, message domain.Message) error {
	return nil
}

type ServiceTestSuite struct {
	suite.Suite
	ctx            context.Context
	pgContainer    *testhelpers.PostgresContainer
	db             interfaces.Execer
	projectService interfaces.ProjectService
	taskService    interfaces.TaskService
	publicService  interfaces.PublicService
	fixtures       *service_fixtures.Fixtures
}

var USER_ONE, USER_TWO, USER_THREE string

func (suite *ServiceTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	suite.pgContainer, err = testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.db, err = infra.NewDatabase("postgres", suite.pgContainer.ConnectionString)
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
	redis := infra.NewInMemory(redisClient)
	rateLimiter := infra.NewRateLimiter(redisClient, 5, 3)
	store := NewStorage(suite.db, redis, rateLimiter)

	notifier := &mockNotifier{}

	suite.projectService = NewProjectService(store, notifier)
	suite.taskService = NewTaskService(store, notifier)
	suite.publicService = NewPublicService(store)

	suite.fixtures = service_fixtures.New(suite.ctx, store)

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
	_, err := suite.db.ExecContext(suite.ctx, "TRUNCATE users CASCADE")
	if err != nil {
		log.Fatal(err)
	}
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *ServiceTestSuite) Cleanup() {
	_, err := suite.db.ExecContext(suite.ctx, "DELETE FROM projects")
	suite.Require().NoError(err)
}

func TestServices(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
