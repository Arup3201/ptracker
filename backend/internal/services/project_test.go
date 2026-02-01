package services

import (
	"testing"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/testhelpers/service_fixtures"
)

func (suite *ServiceTestSuite) TestCreateProject() {
	t := suite.T()

	t.Run("should create project with non-empty id", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Require().NoError(err)
		suite.Require().NotEqual("", id)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should create project with correct values", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Require().NoError(err)
		var p domain.Project
		suite.db.QueryRow(
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should create project without description", func(t *testing.T) {
		sample_name := "Test Project"
		sample_skills := "C++, Java"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			nil,
			&sample_skills,
			USER_ONE,
		)

		suite.Require().NoError(err)
		var p domain.Project
		suite.db.QueryRow(
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal((*string)(nil), p.Description)
		suite.Require().Equal(sample_skills, *p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should create project without skills", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			nil,
			USER_ONE,
		)

		suite.Require().NoError(err)
		var p domain.Project
		suite.db.QueryRow(
			"SELECT "+
				"id, name, description, skills, owner, created_at, updated_at "+
				"FROM projects "+
				"WHERE id=($1)",
			id,
		).Scan(&p.Id, &p.Name, &p.Description, &p.Skills, &p.Owner, &p.CreatedAt, &p.UpdatedAt)
		suite.Require().Equal(sample_name, p.Name)
		suite.Require().Equal(sample_description, *p.Description)
		suite.Require().Equal((*string)(nil), p.Skills)
		suite.Require().Equal(USER_ONE, p.Owner)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should fail to create project with empty name", func(t *testing.T) {
		sample_name := ""
		sample_description := "Test Description"
		sample_skills := "C++, Java"

		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		suite.Require().EqualError(err, "invalid value")

		service.store.Project().Delete(suite.ctx, id)
	})
}

func (suite *ServiceTestSuite) TestGetPrivateProject() {
	t := suite.T()

	t.Run("should get project with correct id", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := service.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		suite.Require().Equal(id, project.Id)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should get project with correct name, description and skills", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := service.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		suite.Require().Equal(sample_name, project.Name)
		suite.Require().Equal(sample_description, *project.Description)
		suite.Require().Equal(sample_skills, *project.Skills)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should get project with role owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := service.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		suite.Require().Equal(domain.ROLE_OWNER, project.Role)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should get project with correct owner", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		project, err := service.GetPrivateProject(suite.ctx, id, USER_ONE)

		suite.Require().NoError(err)
		var u domain.Member
		suite.db.QueryRow(
			"SELECT "+
				"id, username, display_name, email, avatar_url, is_active, created_at, updated_at  "+
				"FROM users "+
				"WHERE id=($1)",
			USER_ONE,
		).Scan(&u.Id, &u.Username, &u.DisplayName, &u.Email, &u.AvatarURL, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		suite.Require().Equal(u.Id, project.Owner.Id)
		suite.Require().Equal(u.Username, project.Owner.Username)
		suite.Require().Equal(u.DisplayName, project.Owner.DisplayName)
		suite.Require().Equal(u.Email, project.Owner.Email)
		suite.Require().Equal(u.AvatarURL, project.Owner.AvatarURL)
		suite.Require().Equal(u.IsActive, project.Owner.IsActive)

		service.store.Project().Delete(suite.ctx, id)
	})
	t.Run("should get forbidden error", func(t *testing.T) {
		sample_name := "Test Project"
		sample_description := "Test Description"
		sample_skills := "C++, Java"
		service := NewProjectService(suite.store)
		id, err := service.CreateProject(suite.ctx,
			sample_name,
			&sample_description,
			&sample_skills,
			USER_ONE,
		)

		_, err = service.GetPrivateProject(suite.ctx, id, USER_TWO)

		suite.Require().EqualError(err, "forbidden")

		service.store.Project().Delete(suite.ctx, id)
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

		service := NewProjectService(suite.store)
		members, err := service.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Require().NoError(err)
		suite.Require().Equal(1, len(members))

		service.store.Project().Delete(suite.ctx, p)
	})
	t.Run("should get 2 members", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)
		suite.fixtures.Role(p, USER_THREE, domain.ROLE_MEMBER)

		service := NewProjectService(suite.store)
		members, err := service.GetProjectMembers(suite.ctx, p, USER_TWO)

		suite.Require().NoError(err)
		suite.Require().Equal(2, len(members))

		service.store.Project().Delete(suite.ctx, p)
	})
	t.Run("should get forbidden error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

		service := NewProjectService(suite.store)
		_, err := service.GetProjectMembers(suite.ctx, p, USER_THREE)

		suite.Require().EqualError(err, "forbidden")

		service.store.Project().Delete(suite.ctx, p)
	})
}

func (suite *ServiceTestSuite) TestRespondToJoinRequests() {
	t := suite.T()

	t.Run("should accept join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		err := service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().NoError(err)
	})
	t.Run("should accept join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)
		suite.Require().Equal(domain.JOIN_STATUS_ACCEPTED, status)
	})
	t.Run("should add role member after accepting join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		var role string
		suite.db.QueryRow(
			"SELECT "+
				"role "+
				"FROM roles "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&role)
		suite.Require().Equal("Member", role)

		suite.Cleanup()
	})
	t.Run("should reject join request", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)
		suite.Require().Equal(domain.JOIN_STATUS_REJECTED, status)
	})
	t.Run("should reject join request without creating role", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		err := suite.db.QueryRow(
			"SELECT "+
				"role "+
				"FROM roles "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan()
		suite.Require().ErrorContains(err, "no rows in result set")

		suite.Cleanup()
	})
	t.Run("should not transition join status from accepted to pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)
		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		err := service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_PENDING)

		suite.Require().EqualError(err, "invalid value")
	})
	t.Run("should not transition join status from rejected to accepted", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)
		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		err := service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().EqualError(err, "invalid value")
	})
	t.Run("should transition join status from rejected to pending", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)
		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_REJECTED)

		service.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_PENDING)

		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)
		suite.Require().Equal(domain.JOIN_STATUS_PENDING, status)
	})
	t.Run("should not change join status from pending due to transaction error", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		// just for testing transaction working fine...
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)
		projectService := NewProjectService(suite.store)

		err := projectService.RespondToJoinRequests(suite.ctx, p, USER_ONE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().ErrorContains(err, "transaction: store role create")
		var status string
		suite.db.QueryRow(
			"SELECT "+
				"status "+
				"FROM join_requests "+
				"WHERE project_id=($1) AND user_id=($2)",
			p, USER_TWO,
		).Scan(&status)
		suite.Require().Equal(status, domain.JOIN_STATUS_PENDING)

		suite.Cleanup()
	})
	t.Run("should be forbidden to change the join status by non-member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		err := service.RespondToJoinRequests(suite.ctx, p, USER_THREE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().ErrorContains(err, "forbidden")
	})
	t.Run("should be forbidden to change the join status by member", func(t *testing.T) {
		p := suite.fixtures.Project(service_fixtures.ProjectParams{
			Title:   "Project Fixture A",
			OwnerID: USER_ONE,
		})
		suite.fixtures.Role(p, USER_THREE, domain.ROLE_MEMBER)
		NewPublicService(suite.store).JoinProject(suite.ctx, p, USER_TWO)
		service := NewProjectService(suite.store)

		err := service.RespondToJoinRequests(suite.ctx, p, USER_THREE, USER_TWO, domain.JOIN_STATUS_ACCEPTED)

		suite.Require().ErrorContains(err, "forbidden")
	})
}
