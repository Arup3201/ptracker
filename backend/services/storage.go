package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/ptracker/interfaces"
	"github.com/ptracker/repositories"
)

type Storage struct {
	mu sync.Mutex

	db repositories.Execer

	sessionRepo     interfaces.SessionRepository
	userRepo        interfaces.UserRepository
	projectRepo     interfaces.ProjectRepository
	taskRepo        interfaces.TaskRepository
	roleRepo        interfaces.RoleRepository
	listRepo        interfaces.ListRepository
	joinRequestRepo interfaces.JoinRequestRepository
	publicRepo      interfaces.PublicRepository

	inMemory interfaces.InMemory
}

func NewStorage(db repositories.Execer,
	memory interfaces.InMemory) interfaces.Store {
	s := &Storage{}
	s.db = db
	s.sessionRepo = repositories.NewSessionRepo(db)
	s.userRepo = repositories.NewUserRepo(db)
	s.projectRepo = repositories.NewProjectRepo(db)
	s.taskRepo = repositories.NewTaskRepo(db)
	s.roleRepo = repositories.NewRoleRepo(db)
	s.listRepo = repositories.NewListRepo(db)
	s.joinRequestRepo = repositories.NewJoinRequestRepo(db)
	s.publicRepo = repositories.NewPublicRepo(db)

	s.inMemory = memory

	return s
}

func (s *Storage) WithTx(ctx context.Context, fn func(txStore interfaces.Store) error) error {
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

func (s *Storage) Session() interfaces.SessionRepository {
	return s.sessionRepo
}

func (s *Storage) User() interfaces.UserRepository {
	return s.userRepo
}

func (s *Storage) Project() interfaces.ProjectRepository {
	return s.projectRepo
}

func (s *Storage) Task() interfaces.TaskRepository {
	return s.taskRepo
}

func (s *Storage) Role() interfaces.RoleRepository {
	return s.roleRepo
}

func (s *Storage) List() interfaces.ListRepository {
	return s.listRepo
}

func (s *Storage) JoinRequest() interfaces.JoinRequestRepository {
	return s.joinRequestRepo
}

func (s *Storage) Public() interfaces.PublicRepository {
	return s.publicRepo
}

func (s *Storage) clone(tx *sql.Tx) interfaces.Store {
	return NewStorage(tx, s.inMemory)
}

func (s *Storage) InMemory() interfaces.InMemory {
	return s.inMemory
}
