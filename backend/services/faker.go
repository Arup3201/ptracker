package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/domain"
)

type fakeStore struct{}

func (f *fakeStore) Project() ProjectRepository {
	return NewFakeProjectRepo()
}

func (f *fakeStore) Role() RoleRepository {
	return NewFakeRoleRepo()
}

func (f *fakeStore) WithTx(ctx context.Context, fn func(store Store) error) error {
	return fn(f)
}

type fakeProjectRepo struct {
	projects map[string]*domain.Project
}

func NewFakeProjectRepo() *fakeProjectRepo {
	return &fakeProjectRepo{
		projects: make(map[string]*domain.Project),
	}
}

func (f *fakeProjectRepo) Create(ctx context.Context,
	name string,
	description, skills *string,
	owner string) (string, error) {
	id := uuid.NewString()
	now := time.Now()
	f.projects[id] = &domain.Project{
		Id:          id,
		Name:        name,
		Description: description,
		Skills:      skills,
		Owner:       owner,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	return id, nil
}

func (f *fakeProjectRepo) All(ctx context.Context, userId string) ([]domain.ProjectSummary, error) {
	projects := []domain.ProjectSummary{}
	for _, p := range f.projects {
		// include project when the user is the owner or otherwise (fake repo doesn't track roles)
		if p.Owner != userId {
			continue
		}

		role := domain.ROLE_MEMBER
		if p.Owner == userId {
			role = domain.ROLE_OWNER
		}

		projects = append(projects, domain.ProjectSummary{
			Id:              p.Id,
			Name:            p.Name,
			Description:     p.Description,
			Skills:          p.Skills,
			Role:            role,
			UnassignedTasks: 0,
			OngoingTasks:    0,
			CompletedTasks:  0,
			AbandonedTasks:  0,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		})
	}

	return projects, nil
}

type fakeRoleRepo struct {
	roles map[string]*domain.Role
}

func NewFakeRoleRepo() *fakeRoleRepo {
	return &fakeRoleRepo{
		roles: make(map[string]*domain.Role),
	}
}

func (f *fakeRoleRepo) Create(ctx context.Context,
	projectId, userId, role string) error {
	now := time.Now()
	f.roles[projectId+userId] = &domain.Role{
		ProjectId: projectId,
		UserId:    userId,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	return nil
}

type fakeUserRepo struct {
	users map[string]*domain.User
}

func NewFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (f *fakeUserRepo) Create(ctx context.Context, idpSubject, idpProvider, username string,
	displayName *string,
	email string,
	avatarUrl *string) (string, error) {
	id := uuid.NewString()
	now := time.Now()
	f.users[id] = &domain.User{
		Id:          id,
		IDPSubject:  idpProvider,
		IDPProvider: idpProvider,
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		AvaterURL:   avatarUrl,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	return id, nil
}
