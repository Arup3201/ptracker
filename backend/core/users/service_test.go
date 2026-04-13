package users

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/core"
	"github.com/ptracker/models"
	"github.com/ptracker/testdata"
	"github.com/ptracker/testhelpers"
	"github.com/ptracker/testhelpers/fixtures"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	pgContainer *testhelpers.PostgresContainer
	db          *gorm.DB
	fixtures    *fixtures.Fixtures
	service     *UserService
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(userServiceTestSuite))
}

func (suite *userServiceTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	suite.pgContainer, err = testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.db, err = gorm.Open(postgres.Open(suite.pgContainer.ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = testdata.TestMigrate(suite.db)
	if err != nil {
		log.Fatal(err)
	}

	repo := NewUserRepository(suite.db)
	suite.service = NewUserService(repo)

	suite.fixtures = fixtures.New(suite.ctx, suite.db)
}

func (suite *userServiceTestSuite) Cleanup() {
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM users").Error
	suite.Require().NoError(err)
}

func (suite *userServiceTestSuite) TestCreate() {
	t := suite.T()

	t.Run("should create user", func(t *testing.T) {
		username := "alice"
		email := "alice@test.com"

		_, err := suite.service.Create(suite.ctx, username, email, nil, nil)

		suite.Require().NoError(err)
	})
	t.Run("should create user with displayname and avatar", func(t *testing.T) {
		dn := "Alice"
		au := "http://avatar"
		username := "alice2"
		email := "alice2@test.com"

		id, err := suite.service.Create(suite.ctx, username, email, &dn, &au)

		suite.Require().NoError(err)

		u, err := gorm.G[models.User](suite.db).Where("id = ?", id).First(suite.ctx)
		suite.Require().NoError(err)
		suite.Require().Equal(dn, *u.DisplayName)
		suite.Require().Equal(au, *u.AvatarURL)
	})
	t.Run("should get invalid value error with empty email", func(t *testing.T) {
		username := "alice"
		email := ""

		_, err := suite.service.Create(suite.ctx, username, email, nil, nil)

		suite.Require().ErrorIs(err, core.ErrInvalidValue)
	})
	t.Run("should get invalid value error with invalid email address", func(t *testing.T) {
		username := "alice"
		email := "alice.test@com"

		_, err := suite.service.Create(suite.ctx, username, email, nil, nil)

		suite.Require().ErrorIs(err, core.ErrInvalidValue)
	})
	t.Run("should get invalid value error with empty username", func(t *testing.T) {
		username := ""
		email := "alice@test.com"

		_, err := suite.service.Create(suite.ctx, username, email, nil, nil)

		suite.Require().ErrorIs(err, core.ErrInvalidValue)
	})
	suite.Cleanup()
}

func (suite *userServiceTestSuite) TestGet() {
	t := suite.T()

	t.Run("should get existing user", func(t *testing.T) {
		username := "alice-service"
		email := "alice@svc.test"
		dn := "Alice"
		au := "http://avatar"
		id := suite.fixtures.InsertUser(models.User{
			Username:    username,
			Email:       email,
			DisplayName: &dn,
			AvatarURL:   &au,
		})

		u, err := suite.service.Get(suite.ctx, id)

		suite.Require().NoError(err)
		suite.Require().Equal(id, u.ID)
		suite.Require().Equal("alice-service", u.Username)
		suite.Require().Equal("alice@svc.test", u.Email)
		suite.Require().NotNil(u.DisplayName)
		suite.Require().Equal(dn, *u.DisplayName)
		suite.Require().NotNil(u.AvatarURL)
		suite.Require().Equal(au, *u.AvatarURL)
	})

	t.Run("should return not found for invalid id", func(t *testing.T) {
		_, err := suite.service.Get(suite.ctx, "invalid-id")

		suite.Require().ErrorIs(err, core.ErrNotFound)
	})
	suite.Cleanup()
}
