package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories"
)

type storage struct {
	mu sync.Mutex

	db interfaces.Execer

	sessionRepo     interfaces.SessionRepository
	userRepo        interfaces.UserRepository
	projectRepo     interfaces.ProjectRepository
	taskRepo        interfaces.TaskRepository
	roleRepo        interfaces.RoleRepository
	assigneeRepo    interfaces.AssigneeRepository
	listRepo        interfaces.ListRepository
	joinRequestRepo interfaces.JoinRequestRepository
	publicRepo      interfaces.PublicRepository

	inMemory    interfaces.InMemory
	rateLimiter interfaces.RateLimiter
}

func NewStorage(db interfaces.Execer,
	memory interfaces.InMemory,
	rateLimiter interfaces.RateLimiter) interfaces.Store {
	s := &storage{}
	s.db = db
	s.sessionRepo = repositories.NewSessionRepo(db)
	s.userRepo = repositories.NewUserRepo(db)
	s.projectRepo = repositories.NewProjectRepo(db)
	s.taskRepo = repositories.NewTaskRepo(db)
	s.roleRepo = repositories.NewRoleRepo(db)
	s.assigneeRepo = repositories.NewAssigneeRepo(db)
	s.listRepo = repositories.NewListRepo(db)
	s.joinRequestRepo = repositories.NewJoinRequestRepo(db)
	s.publicRepo = repositories.NewPublicRepo(db)

	s.inMemory = memory
	s.rateLimiter = rateLimiter

	return s
}

func (s *storage) WithTx(ctx context.Context, fn func(txStore interfaces.Store) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sqlDB, ok := s.db.(*sql.DB)
	if !ok {
		return errors.New("WithTx called on non-root store")
	}

	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure rollback on panic or error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	storeTx := s.clone(tx)

	if err := fn(storeTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *storage) Session() interfaces.SessionRepository {
	return s.sessionRepo
}

func (s *storage) User() interfaces.UserRepository {
	return s.userRepo
}

func (s *storage) Project() interfaces.ProjectRepository {
	return s.projectRepo
}

func (s *storage) Task() interfaces.TaskRepository {
	return s.taskRepo
}

func (s *storage) Role() interfaces.RoleRepository {
	return s.roleRepo
}

func (s *storage) Assignee() interfaces.AssigneeRepository {
	return s.assigneeRepo
}

func (s *storage) List() interfaces.ListRepository {
	return s.listRepo
}

func (s *storage) JoinRequest() interfaces.JoinRequestRepository {
	return s.joinRequestRepo
}

func (s *storage) Public() interfaces.PublicRepository {
	return s.publicRepo
}

func (s *storage) clone(tx *sql.Tx) interfaces.Store {
	return NewStorage(tx, s.inMemory, s.rateLimiter)
}

func (s *storage) InMemory() interfaces.InMemory {
	return s.inMemory
}

func (s *storage) RateLimiter() interfaces.RateLimiter {
	return s.rateLimiter
}
