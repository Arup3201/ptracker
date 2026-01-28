package fake

import (
	"context"
	"sync"

	"github.com/ptracker/stores"
)

type Store struct {
	mu sync.Mutex

	userRepo    *UserRepo
	projectRepo *ProjectRepo
	roleRepo    *RoleRepo
}

func NewStore() *Store {
	s := &Store{}
	s.userRepo = NewUserRepo()
	s.projectRepo = NewProjectRepo()
	s.roleRepo = NewRoleRepo()
	return s
}

func (s *Store) WithTx(ctx context.Context, fn func(txStore stores.Store) error) error {
	// No real transaction — but we preserve semantics
	s.mu.Lock()
	defer s.mu.Unlock()

	return fn(s)
}

func (s *Store) User() stores.UserRepository {
	return s.userRepo
}

func (s *Store) Project() stores.ProjectRepository {
	return s.projectRepo
}

func (s *Store) Role() stores.RoleRepository {
	return s.roleRepo
}
