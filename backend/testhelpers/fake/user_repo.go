package fake

import (
	"context"
	"sync"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, apierr.ErrResourceNotFound
	}

	cp := *u
	return &cp, nil
}

func (r *UserRepo) WithUser(u domain.User) *UserRepo {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := u
	r.users[u.Id] = &cp
	return r
}
