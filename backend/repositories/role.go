package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
)

type RoleRepo struct {
	DB Execer
}

func NewRoleRepo(db Execer) stores.RoleRepository {
	return &RoleRepo{
		DB: db,
	}
}

func (r *RoleRepo) Create(ctx context.Context,
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

func (r *RoleRepo) Get(ctx context.Context, projectId, userId string) (string, error) {
	var role string
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT "+
			"role "+
			"FROM roles "+
			"WHERE project_id=($1) AND user_id=($2)",
		projectId, userId).
		Scan(&role)
	if err != nil {
		return "", fmt.Errorf("postgres query user role: %w", err)
	}

	return role, nil
}

func (r *RoleRepo) CountMembers(ctx context.Context, projectId string) (int, error) {
	var count int
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT "+
			"COUNT(user_id) "+
			"FROM roles "+
			"WHERE project_id=($1) AND role!=($2)",
		projectId, domain.ROLE_OWNER,
	).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("postgres query total members: %w", err)
	}

	return count, nil
}
