package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/ptracker/interfaces"
	"github.com/ptracker/repositories"
)

type SQLStore struct {
	mu sync.Mutex

	db              repositories.Execer
	userRepo        interfaces.UserRepository
	projectRepo     interfaces.ProjectRepository
	roleRepo        interfaces.RoleRepository
	listRepo        interfaces.ListRepository
	joinRequestRepo interfaces.JoinRequestRepository
	publicRepo      interfaces.PublicRepository
}

func NewSQLStore(db repositories.Execer) interfaces.Store {
	s := &SQLStore{}
	s.db = db
	s.userRepo = repositories.NewUserRepo(db)
	s.projectRepo = repositories.NewProjectRepo(db)
	s.roleRepo = repositories.NewRoleRepo(db)
	s.listRepo = repositories.NewListRepo(db)
	s.joinRequestRepo = repositories.NewJoinRequestRepo(db)
	s.publicRepo = repositories.NewPublicRepo(db)

	return s
}

func (s *SQLStore) WithTx(ctx context.Context, fn func(txStore interfaces.Store) error) error {
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

func (s *SQLStore) User() interfaces.UserRepository {
	return s.userRepo
}

func (s *SQLStore) Project() interfaces.ProjectRepository {
	return s.projectRepo
}

func (s *SQLStore) Role() interfaces.RoleRepository {
	return s.roleRepo
}

func (s *SQLStore) List() interfaces.ListRepository {
	return s.listRepo
}

func (s *SQLStore) JoinRequest() interfaces.JoinRequestRepository {
	return s.joinRequestRepo
}

func (s *SQLStore) Public() interfaces.PublicRepository {
	return s.publicRepo
}

func (s *SQLStore) clone(tx *sql.Tx) interfaces.Store {
	return NewSQLStore(tx)
}
