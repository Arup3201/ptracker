package fake

import (
	"context"
	"sync"

	"github.com/ptracker/apierr"
)

type RoleRepo struct {
	mu    sync.RWMutex
	roles map[string]map[string]string // projectId -> userId -> role
}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{
		roles: make(map[string]map[string]string),
	}
}

func (r *RoleRepo) Create(ctx context.Context, projectId, userId, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.roles[projectId]; !ok {
		r.roles[projectId] = make(map[string]string)
	}

	r.roles[projectId][userId] = role
	return nil
}

func (r *RoleRepo) Get(ctx context.Context, projectId, userId string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, ok := r.roles[projectId]
	if !ok {
		return "", apierr.ErrResourceNotFound
	}

	role, ok := users[userId]
	if !ok {
		return "", apierr.ErrResourceNotFound
	}

	return role, nil
}

func (r *RoleRepo) CountMembers(ctx context.Context, projectId string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users, ok := r.roles[projectId]
	if !ok {
		return 0, nil
	}

	return len(users), nil
}

func (r *RoleRepo) WithRole(projectId, userId, role string) *RoleRepo {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.roles[projectId]; !ok {
		r.roles[projectId] = make(map[string]string)
	}

	r.roles[projectId][userId] = role
	return r
}
