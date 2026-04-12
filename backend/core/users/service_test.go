package users

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/core"
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
	repo        *UserRepository
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

	suite.repo = NewUserRepository(suite.db)
	suite.service = NewUserService(suite.repo)

	suite.fixtures = fixtures.New(suite.ctx, suite.db)
}

func (suite *userServiceTestSuite) Cleanup() {
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM projects").Error
	suite.Require().NoError(err)
}

func (suite *userServiceTestSuite) TestGet() {
	t := suite.T()

	t.Run("should get existing user", func(t *testing.T) {
		// create user with display name and avatar via repository
		dn := "Alice"
		au := "http://avatar"
		id, err := suite.repo.Create(suite.ctx, "alice-service", "alice@svc.test", &dn, &au)
		if err != nil {
			t.Fatalf("failed to create user fixture: %v", err)
		}

		u, err := suite.service.Get(suite.ctx, id)

		suite.Cleanup()

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

		suite.Cleanup()

		suite.Require().ErrorIs(err, core.ErrNotFound)
	})
}
