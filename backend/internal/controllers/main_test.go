package controllers

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/infra"
	"github.com/ptracker/internal/services"
	"github.com/ptracker/internal/testdata"
	"github.com/ptracker/internal/testhelpers"
	"github.com/ptracker/internal/testhelpers/controller_fixtures"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db          *gorm.DB
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

	conn, err := infra.NewDatabase(pgContainer.ConnectionString, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	testdata.TestMigrate(conn)

	suite.db = conn

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

	taskService := services.NewTaskService(store, notifier)
	taskController := NewTaskController(taskService)

	handler := http.NewServeMux()
	handler.Handle("POST /projects", HTTPErrorHandler(projectController.Create))
	handler.Handle("GET /projects", HTTPErrorHandler(projectController.List))
	handler.Handle("GET /projects/{id}/members", HTTPErrorHandler(projectController.ListMembers))
	handler.Handle("PUT /projects/{project_id}/tasks/{task_id}", HTTPErrorHandler(taskController.Update))

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
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM projects").Error
	suite.Require().NoError(err)
}

func TestControllers(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}
