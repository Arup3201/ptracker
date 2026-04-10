package services

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/repositories/models"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
	"gorm.io/gorm"
)

func (suite *ServiceTestSuite) TestCreateTask() {
	t := suite.T()

	t.Run("should create task", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		sample_description := "sample description"

		_, _, err := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{})

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create task with unassigned status", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		sample_description := "sample description"

		taskId, _, _ := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{})
		var status string
		suite.db.WithContext(suite.ctx).
			Table("tasks").
			Select("status").
			Where("id = ?", taskId).
			Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(domain.TASK_STATUS_UNASSIGNED, status)
	})
	t.Run("should be forbidden for user two", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		sample_title := "sample task"
		sample_description := "sample description"

		_, _, err := suite.taskService.CreateTask(suite.ctx,
			p, USER_TWO,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{})

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be invalid with empty task title", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := ""
		sample_description := "sample description"

		_, _, err := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{})

		suite.Cleanup()

		suite.Require().ErrorContains(err, "invalid value")
	})
	t.Run("should create tasks with assignees", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)
		sample_title := "sample task"
		sample_description := "sample description"

		_, _, err := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{USER_TWO, USER_THREE})

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should create tasks with 2 assignees", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)
		sample_title := "sample task"
		sample_description := "sample description"

		_, warnings, _ := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{USER_TWO, USER_THREE})

		suite.Cleanup()

		suite.Require().Equal(0, len(warnings))
	})
	t.Run("should create task with correct assignee values", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)
		sample_title := "sample task"
		sample_description := "sample description"

		id, _, _ := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{USER_TWO, USER_THREE})

		rows, _ := gorm.G[models.Assignee](suite.db).Where("project_id=? AND task_id=?", p, id).Find(suite.ctx)
		suite.Cleanup()

		suite.Require().Equal(2, len(rows))
	})
	t.Run("should give one warning with 1 invalid assignee", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)
		sample_title := "sample task"
		sample_description := "sample description"

		_, warnings, _ := suite.taskService.CreateTask(suite.ctx,
			p, USER_ONE,
			sample_title, sample_description, domain.TASK_STATUS_UNASSIGNED,
			[]string{USER_TWO, "asdfd"})

		suite.Cleanup()

		suite.Require().Equal(1, len(warnings))
	})
}

func (suite *ServiceTestSuite) TestListTasks() {
	t := suite.T()

	t.Run("should get list of tasks", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		_, err := suite.taskService.ListTasks(suite.ctx, p, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get list of tasks for member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)

		_, err := suite.taskService.ListTasks(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get forbidden error for non-member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: p,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		_, err := suite.taskService.ListTasks(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}

func (suite *ServiceTestSuite) TestGetTask() {
	t := suite.T()

	t.Run("should get task", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		_, err := suite.taskService.GetTask(suite.ctx, projectId, taskId, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should get task with id title status", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     sample_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		task, _ := suite.taskService.GetTask(suite.ctx, projectId, taskId, USER_ONE)

		suite.Cleanup()

		suite.Require().Equal(taskId, task.ID)
		suite.Require().Equal(sample_title, task.Title)
		suite.Require().Equal(domain.TASK_STATUS_UNASSIGNED, task.Status)
	})
	t.Run("should be forbidden for non-member", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		sample_title := "sample task"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     sample_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		_, err := suite.taskService.GetTask(suite.ctx, projectId, taskId, USER_TWO)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}

func (suite *ServiceTestSuite) TestUpdateTask() {
	t := suite.T()

	t.Run("should update task", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		updatedTaskTitle := "Project title updated"

		err := suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_ONE,
			&updatedTaskTitle, nil, nil,
			nil, nil)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should update task with title", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		updatedTaskTitle := "Project title updated"

		suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_ONE,
			&updatedTaskTitle, nil, nil,
			nil, nil)

		var title string
		suite.db.WithContext(suite.ctx).
			Table("tasks").
			Select("title").
			Where("id = ?", taskId).
			Scan(&title)

		suite.Cleanup()

		suite.Require().Equal(updatedTaskTitle, title)
	})
	t.Run("should update description and status", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		task_title := "Project Task A"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     task_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		updatedTaskDesc := "Project description updated"
		updatedStatus := domain.TASK_STATUS_ABANDONED

		suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_ONE,
			nil, &updatedTaskDesc, &updatedStatus,
			nil, nil)

		task, _ := gorm.G[models.Task](suite.db).Where("id = ?", taskId).First(suite.ctx)

		suite.Cleanup()

		suite.Require().NotNil(task.Description)
		suite.Require().Equal(updatedTaskDesc, *task.Description)
		suite.Require().Equal(domain.TASK_STATUS_ABANDONED, task.Status.String)
	})
	t.Run("should update assignees by adding", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		task_title := "Project Task A"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     task_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})

		suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_ONE,
			nil, nil, nil,
			[]string{USER_TWO}, nil)

		rows, _ := gorm.G[models.Assignee](suite.db).Where("task_id = ?", taskId).Find(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(1, len(rows))
	})
	t.Run("should update task by removing assignees", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		task_title := "Project Task A"
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     task_title,
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
			Assignees: []string{USER_TWO},
		})

		suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_ONE,
			nil, nil, nil,
			nil, []string{USER_TWO})

		rows, _ := gorm.G[models.Assignee](suite.db).Where("task_id = ?", taskId).Find(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(0, len(rows))
	})
	t.Run("should be forbidden to update title as member", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		updatedTaskTitle := "Project title updated"

		err := suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_TWO,
			&updatedTaskTitle, nil, nil,
			nil, nil)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be able to update title as assignee", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
			Assignees: []string{USER_TWO},
		})
		updatedTaskTitle := "Project title updated"

		err := suite.taskService.UpdateTask(suite.ctx,
			projectId, taskId, USER_TWO,
			&updatedTaskTitle, nil, nil,
			nil, nil)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
}

func (suite *ServiceTestSuite) TestAddComment() {
	t := suite.T()

	t.Run("should add comment", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		sampleComment := "Hello there!"

		id, err := suite.taskService.AddComment(suite.ctx,
			projectId, taskId, USER_TWO, sampleComment)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotEqual("", id)
	})
	t.Run("should add comment with text", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		sampleComment := "Hello there!"

		id, _ := suite.taskService.AddComment(suite.ctx,
			projectId, taskId, USER_TWO, sampleComment)

		var comment string
		suite.db.WithContext(suite.ctx).
			Table("comments").
			Select("content").
			Where("id = ?", id).
			Scan(&comment)

		suite.Cleanup()

		suite.Require().Equal(sampleComment, comment)
	})
	t.Run("should be forbidden for non-member", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		sampleComment := "Hello there!"

		_, err := suite.taskService.AddComment(suite.ctx,
			projectId, taskId, USER_TWO, sampleComment)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be invalid with empty comment", func(t *testing.T) {
		projectId := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(projectId, USER_TWO, domain.ROLE_MEMBER)
		taskId := suite.fixtures.Task(service_fixtures.TaskParams{
			Title:     "Project Task A",
			ProjectId: projectId,
			Status:    domain.TASK_STATUS_UNASSIGNED,
		})
		sampleComment := ""

		_, err := suite.taskService.AddComment(suite.ctx,
			projectId, taskId, USER_TWO, sampleComment)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "invalid value")
	})
}
