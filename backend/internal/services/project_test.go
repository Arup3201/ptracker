package services

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/repositories/models"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
	"gorm.io/gorm"
)

func (suite *ServiceTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("should create project with non empty id", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().NotEqual("", id)
	})
	t.Run("should create project with correct values", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, _ := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		p, _ := gorm.G[models.Project](suite.db).
			Where("id = ?", id).First(suite.ctx)
		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.OwnerID)
	})
	t.Run("should create project without description", func(t *testing.T) {
		sample_name := "Test Project"
		sample_skills := "C++, Java"

		id, _ := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			nil,
			&sample_skills,
			USER_ONE,
		)

		p, _ := gorm.G[models.Project](suite.db).
			Where("id = ?", id).First(suite.ctx)
		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal((*string)(nil), p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.OwnerID)

	})
	t.Run("should create project without skills", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			nil,
			USER_ONE,
		)

		suite.Require().NoError(err)

		p, _ := gorm.G[models.Project](suite.db).
			Where("id = ?", id).First(suite.ctx)
		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal((*string)(nil), p.Skills)
		suite.Require().Equal(USER_ONE, p.OwnerID)
	})
	t.Run("should fail to create project with empty name", func(t *testing.T) {
		sample_name := ""
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		_, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Cleanup()

		suite.Require().EqualError(err, "invalid value")
	})
}

func (suite *ServiceTestSuite) TestGetPrivateProject() {
	t := suite.T()

	t.Run("should get project with correct id", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := suite.projectService.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(id, project.ID)
	})
	t.Run("should get project with correct name, description and skills", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := suite.projectService.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(sample_name, project.Name)
		suite.Require().Equal(sample_description, *project.Description)
		suite.Require().Equal(sample_skills, *project.Skills)
	})
	t.Run("should get project with role owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := suite.projectService.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(domain.ROLE_OWNER, project.Role)
	})
	t.Run("should get project with correct owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := suite.projectService.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		user, _ := gorm.G[models.User](suite.db).
			Where("id = ?", USER_ONE).First(suite.ctx)
		suite.Cleanup()

		suite.Require().Equal(user.ID, project.Owner.UserID)
		suite.Require().Equal(user.Username, project.Owner.Username)
		suite.Require().Equal(user.DisplayName, project.Owner.DisplayName)
		suite.Require().Equal(user.Email, project.Owner.Email)
		suite.Require().Equal(user.AvatarURL, project.Owner.AvatarURL)
	})
	t.Run("should get forbidden error", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		id, err := suite.projectService.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		_, err = suite.projectService.GetPrivateProject(suite.ctx, id, USER_TWO)

		suite.Cleanup()

		suite.Require().EqualError(err, "forbidden")
	})
}

func (suite *ServiceTestSuite) TestGetMembers() {
	t := suite.T()

	t.Run("should get 1 member and owner", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)

		members, err := suite.projectService.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(members))
	})
	t.Run("should get 2 members and owner", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)

		members, err := suite.projectService.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(3, len(members))
	})
	t.Run("should get forbidden error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)

		_, err := suite.projectService.GetProjectMembers(suite.ctx, p, USER_THREE)

		suite.Cleanup()

		suite.Require().EqualError(err, "forbidden")
	})
}

func (suite *ServiceTestSuite) TestRespondToJoinRequests() {
	t := suite.T()

	t.Run("should accept join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Cleanup()

		suite.Require().NoError(err)
	})
	t.Run("should accept join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		joinRequest, _ := gorm.G[models.JoinRequest](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_ACCEPTED, joinRequest.Status.String)
	})
	t.Run("should add role member after accepting join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		membership, _ := gorm.G[models.Membership](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal("Member", membership.Role.String)
	})
	t.Run("should reject join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		joinRequest, _ := gorm.G[models.JoinRequest](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_REJECTED, joinRequest.Status.String)
	})
	t.Run("should reject join request without creating role", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		_, err := gorm.G[models.Membership](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Error(err, gorm.ErrRecordNotFound)
	})
	t.Run("should not transition join status from accepted to pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_PENDING)

		suite.Cleanup()

		suite.Require().EqualError(err, "invalid value")
	})
	t.Run("should not transition join status from rejected to accepted", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Cleanup()

		suite.Require().EqualError(err, "invalid value")
	})
	t.Run("should transition join status from rejected to pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_PENDING)

		joinRequest, _ := gorm.G[models.JoinRequest](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_PENDING, joinRequest.Status.String)
	})
	t.Run("should not change join status from pending due to transaction error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)
		// just for testing transaction working fine...
		suite.fixtures.Membership(p, USER_TWO, domain.ROLE_MEMBER)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().ErrorContains(err, "transaction: store role create")
		joinRequest, _ := gorm.G[models.JoinRequest](suite.db).
			Where("project_id = ? AND user_id = ?", p, USER_TWO).First(suite.ctx)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_PENDING, joinRequest.Status.String)
	})
	t.Run("should be forbidden to change the join status by non-member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_THREE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be forbidden to change the join status by member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Membership(p, USER_THREE, domain.ROLE_MEMBER)
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_THREE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}
