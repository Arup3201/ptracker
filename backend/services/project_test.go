package services

import (
	"testing"

	"github.com/ptracker/domain"
	"github.com/ptracker/testhelpers/service_fixtures"
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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
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
		suite.fixtures.Role(p, USER_ONE, domain.ROLE_OWNER)
		suite.fixtures.Role(p, USER_TWO, domain.ROLE_MEMBER)

		service := NewProjectService(suite.store)
		_, err := service.GetProjectMembers(suite.ctx, p, USER_THREE)

		suite.Require().EqualError(err, "forbidden")

		service.store.Project().Delete(suite.ctx, p)
	})
}
