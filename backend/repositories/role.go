package repositories

import (
	"context"
	"fmt"
	"time"
)

type roleRepo struct {
	DB Execer
}

func NewRoleRepo(db Execer) *roleRepo {
	return &roleRepo{
		DB: db,
	}
}

func (r *roleRepo) Create(ctx context.Context,
	projectId, userId, role string) error {
	now := time.Now()

	_, err := r.DB.ExecContext(ctx, "INSERT INTO "+
		"roles(project_id, user_id, role, created_at, updated_at) "+
		"VALUES($1, $2, $3, $4, $4)",
		projectId, userId, role, now)
	if err != nil {
		return fmt.Errorf("db exec context: %w", err)
	}

	return nil
}
