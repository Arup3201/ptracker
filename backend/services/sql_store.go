package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/ptracker/repositories"
	"github.com/ptracker/stores"
)

type SQLStore struct {
	mu sync.Mutex

	db          repositories.Execer
	userRepo    stores.UserRepository
	projectRepo stores.ProjectRepository
	roleRepo    stores.RoleRepository
	listRepo    stores.ListRepository
}

func NewSQLStore(db repositories.Execer) stores.Store {
	s := &SQLStore{}
	s.db = db
	s.userRepo = repositories.NewUserRepo(db)
	s.projectRepo = repositories.NewProjectRepo(db)
	s.roleRepo = repositories.NewRoleRepo(db)
	s.listRepo = repositories.NewListRepo(db)
	return s
}

func (s *SQLStore) WithTx(ctx context.Context, fn func(txStore stores.Store) error) error {
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

func (s *SQLStore) User() stores.UserRepository {
	return s.userRepo
}

func (s *SQLStore) Project() stores.ProjectRepository {
	return s.projectRepo
}

func (s *SQLStore) Role() stores.RoleRepository {
	return s.roleRepo
}

func (s *SQLStore) List() stores.ListRepository {
	return s.listRepo
}

func (s *SQLStore) clone(tx *sql.Tx) stores.Store {
	return NewSQLStore(tx)
}
