package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type projectRepo struct {
	DB Execer
}

func NewProjectRepo(db Execer) *projectRepo {
	return &projectRepo{
		DB: db,
	}
}

func (r *projectRepo) Create(ctx context.Context,
	name string,
	description, skills *string,
	owner string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.DB.ExecContext(ctx, "INSERT INTO "+
		"projects(id, name, description, skills, owner, created_at, updated_at) "+
		"VALUES($1, $2, $3, $4, $5, $6, $6)",
		id, name, description, skills, owner, now)
	if err != nil {
		return "", fmt.Errorf("db exec context: %w", err)
	}

	return id, nil
}
