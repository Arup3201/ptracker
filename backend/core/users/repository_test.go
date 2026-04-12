package users

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/core"
	"github.com/ptracker/core/models"
	"github.com/ptracker/testdata"
	"github.com/ptracker/testhelpers"
	"github.com/ptracker/testhelpers/fixtures"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var USER_ONE string

type userRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *gorm.DB
	fixtures    *fixtures.Fixtures
	repo        *UserRepository
	ctx         context.Context
}

func (suite *userRepositoryTestSuite) SetupSuite() {
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

	suite.repo = NewUserRepository(suite.db)

	err = testdata.TestMigrate(suite.db)
	if err != nil {
		log.Fatal(err)
	}

	suite.fixtures = fixtures.New(suite.ctx, suite.db)

	USER_ONE = suite.fixtures.InsertUser(fixtures.RandomUserRow())
}

func (suite *userRepositoryTestSuite) Cleanup() {
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM projects").Error
	suite.Require().NoError(err)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(userRepositoryTestSuite))
}

func (suite *userRepositoryTestSuite) TestCreate() {
	t := suite.T()

	t.Run("should create user", func(t *testing.T) {
		id, err := suite.repo.Create(suite.ctx, "alice", "alice@test.com", nil, nil)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotEmpty(id)
	})

	t.Run("should create user with display name and avatar url", func(t *testing.T) {
		dn := "Alice"
		au := "http://avatar"
		id, err := suite.repo.Create(suite.ctx, "alice2", "alice2@test.com", &dn, &au)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotEmpty(id)

		u, err := gorm.G[models.User](suite.db).Where("id = ?", id).First(suite.ctx)
		suite.Require().NoError(err)
		suite.Require().Equal(dn, *u.DisplayName)
		suite.Require().Equal(au, *u.AvatarURL)
	})
}

func (suite *userRepositoryTestSuite) TestGet() {
	t := suite.T()

	t.Run("should get existing user", func(t *testing.T) {
		id := suite.fixtures.InsertUser(fixtures.RandomUserRow())

		u, err := suite.repo.Get(suite.ctx, id)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(id, u.ID)
	})

	t.Run("should return not found for invalid id", func(t *testing.T) {
		_, err := suite.repo.Get(suite.ctx, "invalid")

		suite.Cleanup()

		suite.Require().ErrorIs(err, core.ErrNotFound)
	})
}
