package projects

import (
	"context"
	"log"
	"testing"

	"github.com/ptracker/core"
	"github.com/ptracker/core/models"
	"github.com/ptracker/testdata"
	"github.com/ptracker/testhelpers"
	"github.com/ptracker/testhelpers/repo_fixtures"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var USER_ONE, USER_TWO, USER_THREE string

type projectRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *gorm.DB
	fixtures    *repo_fixtures.Fixtures
	ctx         context.Context
}

func (suite *projectRepositoryTestSuite) SetupSuite() {
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

	suite.fixtures = repo_fixtures.New(suite.ctx, suite.db)

	USER_ONE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
	USER_TWO = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
	USER_THREE = suite.fixtures.InsertUser(repo_fixtures.RandomUserRow())
}

func (suite *projectRepositoryTestSuite) Cleanup() {
	err := suite.db.WithContext(suite.ctx).
		Exec("DELETE FROM projects").Error
	suite.Require().NoError(err)
}

func TestProjectRepository(t *testing.T) {
	suite.Run(t, new(projectRepositoryTestSuite))
}

func (suite *projectRepositoryTestSuite) TestProjectCreate() {
	t := suite.T()

	t.Run("should create project with title description and skills", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		repo := NewProjectRepository(suite.db)
		_, err := repo.Create(suite.ctx,
			sample_name, &sample_description, &sample_skills,
			USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create project with only title", func(t *testing.T) {
		sample_name := "Test Project"

		repo := NewProjectRepository(suite.db)
		_, err := repo.Create(suite.ctx,
			sample_name, nil, nil,
			USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create project with title and description", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		repo := NewProjectRepository(suite.db)
		_, err := repo.Create(suite.ctx,
			sample_name, &sample_description, &sample_skills,
			USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should not create project with invalid user", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		repo := NewProjectRepository(suite.db)
		_, err := repo.Create(suite.ctx,
			sample_name, &sample_description, &sample_skills,
			USER_ONE+"1234")

		suite.Cleanup()

		suite.Require().Error(err)
	})
}

func (suite *projectRepositoryTestSuite) TestProjectGet() {
	t := suite.T()

	t.Run("should get project summary with title description and skills", func(t *testing.T) {
		name, description, skills := "test project", "test description", "c, python"
		projectID := suite.fixtures.InsertProject(models.Project{
			Name:        name,
			Description: &description,
			Skills:      &skills,
			OwnerID:     USER_ONE,
		})
		repo := NewProjectRepository(suite.db)

		project, err := repo.Get(suite.ctx, projectID)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(projectID, project.ID)
		suite.Require().Equal(name, project.Name)
		suite.Require().Equal(description, *project.Description)
		suite.Require().Equal(skills, *project.Skills)
		suite.Require().Equal(USER_ONE, project.OwnerID)
	})
	t.Run("should get project summary with 1 ongoing task", func(t *testing.T) {
		projectID := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(projectID, core.TASK_STATUS_ONGOING))
		repo := NewProjectRepository(suite.db)

		project, _ := repo.Get(suite.ctx, projectID)

		suite.Cleanup()

		suite.Require().EqualValues(1, project.OngoingTasks)
	})
	t.Run("should get project summary with 2 completed and 1 unassigned task", func(t *testing.T) {
		projectID := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(projectID, core.TASK_STATUS_ONGOING))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(projectID, core.TASK_STATUS_COMPLETED))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(projectID, core.TASK_STATUS_COMPLETED))
		repo := NewProjectRepository(suite.db)

		project, _ := repo.Get(suite.ctx, projectID)

		suite.Cleanup()

		suite.Require().EqualValues(1, project.OngoingTasks)
		suite.Require().EqualValues(2, project.CompletedTasks)
	})
	t.Run("should not get any project with invalid id", func(t *testing.T) {
		projectID := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		repo := NewProjectRepository(suite.db)

		_, err := repo.Get(suite.ctx, projectID+"123")

		suite.Cleanup()

		suite.Require().ErrorIs(err, core.ErrNotFound)
	})
}

func (suite *projectRepositoryTestSuite) TestProjectList() {
	t := suite.T()

	t.Run("should return empty list", func(t *testing.T) {
		repo := NewProjectRepository(suite.db)
		projects, err := repo.List(suite.ctx, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(projects)
		suite.Require().Equal(0, len(projects))
	})
	t.Run("should return 2 projects", func(t *testing.T) {
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))

		repo := NewProjectRepository(suite.db)
		projects, err := repo.List(suite.ctx, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(projects))
	})
}

func (suite *projectRepositoryTestSuite) TestProjectRecentlyCreated() {
	t := suite.T()

	t.Run("should get no recently created projects", func(t *testing.T) {
		repo := NewProjectRepository(suite.db)

		projects, err := repo.RecentlyCreated(suite.ctx, USER_ONE, 10)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(projects)
		suite.Require().Equal(0, len(projects))
	})
	t.Run("should get 2 recently created projects", func(t *testing.T) {
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		repo := NewProjectRepository(suite.db)

		projects, err := repo.RecentlyCreated(suite.ctx, USER_ONE, 10)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(projects))
	})
}

func (suite *projectRepositoryTestSuite) TestProjectRecentlyJoined() {
	t := suite.T()

	t.Run("should get no recently joined projects", func(t *testing.T) {
		repo := NewProjectRepository(suite.db)

		projects, err := repo.RecentlyJoined(suite.ctx, USER_TWO, 10)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(projects)
		suite.Require().Equal(0, len(projects))
	})
	t.Run("should get 1 recently joined projects", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertMember(repo_fixtures.GetMemberRow(p, USER_TWO, core.ROLE_MEMBER))
		repo := NewProjectRepository(suite.db)

		projects, err := repo.RecentlyJoined(suite.ctx, USER_TWO, 10)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(1, len(projects))
	})
	t.Run("should get 2 recently joined projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertMember(repo_fixtures.GetMemberRow(p1, USER_TWO, core.ROLE_MEMBER))
		suite.fixtures.InsertMember(repo_fixtures.GetMemberRow(p2, USER_TWO, core.ROLE_MEMBER))
		repo := NewProjectRepository(suite.db)

		projects, err := repo.RecentlyJoined(suite.ctx, USER_TWO, 10)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(projects))
	})
}
