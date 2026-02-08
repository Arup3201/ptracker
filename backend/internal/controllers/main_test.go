package controllers

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/infra"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/services"
	"github.com/ptracker/internal/testhelpers"
	"github.com/ptracker/internal/testhelpers/controller_fixtures"
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

type ControllerTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	fixtures    *controller_fixtures.ControllerFixtures
	db          interfaces.Execer
	ctx         context.Context
}

var (
	USER_ONE, USER_TWO string
)

func (suite *ControllerTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer

	dbConnection, err := infra.NewDatabase("postgres", pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.db = dbConnection

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
	store := services.NewStorage(suite.db, redis, rateLimiter)

	suite.fixtures = controller_fixtures.NewControllerFixtures(suite.ctx, store)

	notifier := &mockNotifier{}
	projectService := services.NewProjectService(store, notifier)
	projectController := NewProjectController(projectService)

	handler := http.NewServeMux()
	handler.Handle("POST /projects", HTTPErrorHandler(projectController.Create))
	handler.Handle("GET /projects", HTTPErrorHandler(projectController.List))
	handler.Handle("GET /projects/{id}/members", HTTPErrorHandler(projectController.ListMembers))

	suite.fixtures.Handler = handler

	USER_ONE = suite.fixtures.User(controller_fixtures.UserParams{
		IDPSubject:  "sub-123",
		IDPProvider: "google",
		Username:    "alice",
		Email:       "alice@example.com",
	})
	USER_TWO = suite.fixtures.User(controller_fixtures.UserParams{
		IDPSubject:  "sub-234",
		IDPProvider: "google",
		Username:    "bob",
		Email:       "bob@example.com",
	})
}

func (suite *ControllerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal(err)
	}
}

func (suite *ControllerTestSuite) Cleanup() {
	_, err := suite.db.ExecContext(suite.ctx, "DELETE FROM projects")
	suite.Require().NoError(err)
}

func TestControllers(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
