package repositories

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/testhelpers/repo_fixtures"
)

func (suite *RepositoryTestSuite) TestProjects() {
	t := suite.T()

	t.Run("should return empty list", func(t *testing.T) {
		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PrivateProjects(suite.ctx, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(projects)
		suite.Require().Equal(0, len(projects))
	})
	t.Run("should return 2 projects", func(t *testing.T) {
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))

		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PrivateProjects(suite.ctx, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		expected := 2
		actual := len(projects)
		suite.Require().Equal(expected, actual)
	})
}

func (suite *RepositoryTestSuite) TestPublicProjects() {
	t := suite.T()

	t.Run("should get empty list", func(t *testing.T) {
		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PublicProjects(suite.ctx, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(projects)
		suite.Require().Equal(0, len(projects))
	})
	t.Run("should get 2 public projects", func(t *testing.T) {
		p1 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		p2 := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))

		listRepo := NewListRepo(suite.db)
		projects, err := listRepo.PublicProjects(suite.ctx, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(projects))
		suite.Require().ElementsMatch(
			[]string{p1, p2},
			[]string{projects[0].Id, projects[1].Id},
		)
	})
}

func (suite *RepositoryTestSuite) TestJoinRequests() {
	t := suite.T()

	t.Run("should get empty list", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))

		listRepo := NewListRepo(suite.db)
		joinRequests, err := listRepo.JoinRequests(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(joinRequests)
		suite.Require().Equal(0, len(joinRequests))
	})
	t.Run("should get list of join requests", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertJoinRequest(repo_fixtures.GetJoinRequest(p, USER_TWO))

		listRepo := NewListRepo(suite.db)
		_, err := listRepo.JoinRequests(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get 2 join requests", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertJoinRequest(repo_fixtures.GetJoinRequest(p, USER_TWO))
		suite.fixtures.InsertJoinRequest(repo_fixtures.GetJoinRequest(p, USER_THREE))

		listRepo := NewListRepo(suite.db)
		joinRequests, _ := listRepo.JoinRequests(suite.ctx, p)

		suite.Cleanup()

		suite.Require().Equal(2, len(joinRequests))
	})
	t.Run("should get join requests from user two and three", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertJoinRequest(repo_fixtures.GetJoinRequest(p, USER_TWO))
		suite.fixtures.InsertJoinRequest(repo_fixtures.GetJoinRequest(p, USER_THREE))

		listRepo := NewListRepo(suite.db)
		joinRequests, _ := listRepo.JoinRequests(suite.ctx, p)

		suite.Cleanup()

		suite.Require().ElementsMatch(
			[]string{joinRequests[0].UserId, joinRequests[1].UserId},
			[]string{USER_TWO, USER_THREE},
		)
	})
}

func (suite *RepositoryTestSuite) TestListTasks() {
	t := suite.T()

	t.Run("should get empty list", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		repo := NewListRepo(suite.db)

		tasks, err := repo.Tasks(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(tasks)
		suite.Require().Equal(0, len(tasks))
	})
	t.Run("should get list of tasks", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		repo := NewListRepo(suite.db)

		_, err := repo.Tasks(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get list of 2 tasks", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		repo := NewListRepo(suite.db)

		tasks, _ := repo.Tasks(suite.ctx, p)

		suite.Cleanup()

		suite.Require().Equal(2, len(tasks))
	})
	t.Run("should get list of 2 tasks with ids", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		t1 := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		t2 := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		repo := NewListRepo(suite.db)

		tasks, _ := repo.Tasks(suite.ctx, p)

		suite.Cleanup()

		suite.Require().ElementsMatch(
			[]string{t1, t2},
			[]string{tasks[0].Id, tasks[1].Id},
		)
	})
	t.Run("should get 1 task with 2 assignees", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_TWO, domain.ROLE_MEMBER))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_THREE, domain.ROLE_MEMBER))
		task := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, task, USER_TWO))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, task, USER_THREE))
		repo := NewListRepo(suite.db)

		tasks, _ := repo.Tasks(suite.ctx, p)

		suite.Cleanup()

		suite.Require().ElementsMatch(
			[]string{USER_TWO, USER_THREE},
			[]string{tasks[0].Assignees[0].UserId, tasks[0].Assignees[1].UserId},
		)
	})
}

func (suite *RepositoryTestSuite) TestListMembers() {
	t := suite.T()

	t.Run("should get empty list", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		repo := NewListRepo(suite.db)

		members, err := repo.Members(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotNil(members)
		suite.Require().Equal(0, len(members))
	})
	t.Run("should list members", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_TWO, domain.ROLE_MEMBER))
		repo := NewListRepo(suite.db)

		members, err := repo.Members(suite.ctx, p)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(1, len(members))
		suite.Require().Equal(USER_TWO, members[0].Id)
	})
}

func (suite *RepositoryTestSuite) TestListAssignees() {
	t := suite.T()

	t.Run("should list assignees", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_TWO, domain.ROLE_MEMBER))
		taskId := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskId, USER_TWO))
		repo := NewListRepo(suite.db)

		_, err := repo.Assignees(suite.ctx, taskId)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should list 2 assignees", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_TWO, domain.ROLE_MEMBER))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_THREE, domain.ROLE_MEMBER))
		taskId := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskId, USER_TWO))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskId, USER_THREE))
		repo := NewListRepo(suite.db)

		assignees, _ := repo.Assignees(suite.ctx, taskId)

		suite.Cleanup()

		suite.Require().Equal(2, len(assignees))
	})
	t.Run("should list 2 assignees with id", func(t *testing.T) {
		p := suite.fixtures.InsertProject(repo_fixtures.RandomProjectRow(USER_ONE))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_TWO, domain.ROLE_MEMBER))
		suite.fixtures.InsertRole(repo_fixtures.GetRoleRow(p, USER_THREE, domain.ROLE_MEMBER))
		taskId := suite.fixtures.InsertTask(repo_fixtures.RandomTaskRow(p, domain.TASK_STATUS_UNASSIGNED))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskId, USER_TWO))
		suite.fixtures.InsertAssignee(repo_fixtures.GetAssigneeRow(p, taskId, USER_THREE))
		repo := NewListRepo(suite.db)

		assignees, _ := repo.Assignees(suite.ctx, taskId)

		suite.Cleanup()

		suite.ElementsMatch(
			[]string{USER_TWO, USER_THREE},
			[]string{assignees[0].UserId, assignees[1].UserId},
		)
	})
}
