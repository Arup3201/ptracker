package services

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
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

		var p domain.Project
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)

		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)
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

		var p domain.Project
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)

		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal((*string)(nil), p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)
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
		var p domain.Project
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)

		suite.Cleanup()

		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal((*string)(nil), p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)
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
		suite.Require().Equal(id, project.Id)
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
		var m domain.Member
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"id, username, display_name, email, avatar_url, is_active, created_at, updated_at  "+
				"FROM users "+
				"WHERE id=($1)",
			USER_ONE,
		).Scan(&m.UserId, &m.Username, &m.DisplayName, &m.Email, &m.AvatarURL, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)

		suite.Cleanup()

		suite.Require().Equal(m.UserId, project.Owner.UserId)
		suite.Require().Equal(m.Username, project.Owner.Username)
		suite.Require().Equal(m.DisplayName, project.Owner.DisplayName)
		suite.Require().Equal(m.Email, project.Owner.Email)
		suite.Require().Equal(m.AvatarURL, project.Owner.AvatarURL)
		suite.Require().Equal(m.IsActive, project.Owner.IsActive)
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

	t.Run("should get 1 member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

		members, err := suite.projectService.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(1, len(members))
	})
	t.Run("should get 2 members", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Role(p, USER_THREE, domain.ROLE_MEMBER)

		members, err := suite.projectService.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Cleanup()

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(members))
	})
	t.Run("should get forbidden error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

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

		var status string
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_ACCEPTED, status)
	})
	t.Run("should add role member after accepting join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		var role string
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"role "+
				"FROM roles "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&role)

		suite.Cleanup()

		suite.Require().Equal("Member", role)
	})
	t.Run("should reject join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		var status string
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_REJECTED, status)
	})
	t.Run("should reject join request without creating role", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		err := suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"role "+
				"FROM roles "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan()
		suite.Cleanup()

		suite.Require().ErrorContains(err, "no rows in result set")
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

		var status string
		suite.db.QueryRowContext(suite.ctx,
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(domain.JOIN_STATUS_PENDING, status)
	})
	t.Run("should not change join status from pending due to transaction error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)
		// just for testing transaction working fine...
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().ErrorContains(err, "transaction: store role create")
		var status string
		suite.db.QueryRowContext(
			suite.ctx,
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)

		suite.Cleanup()

		suite.Require().Equal(status, domain.JOIN_STATUS_PENDING)
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
		suite.fixtures.Role(p, USER_THREE, domain.ROLE_MEMBER)
		suite.publicService.JoinProject(suite.ctx, p, USER_TWO)

		err := suite.projectService.RespondToJoinRequests(suite.ctx, p, USER_THREE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Cleanup()

		suite.Require().ErrorContains(err, "forbidden")
	})
}
