package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
)

type RoleRepo struct {
	db Execer
}

func NewRoleRepo(db Execer) interfaces.RoleRepository {
	return &RoleRepo{
		db: db,
	}
}

func (r *RoleRepo) Create(ctx context.Context,
	projectId, userId, role string) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, "INSERT INTO "+
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
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"role "+
			"FROM roles "+
			"WHERE project_id=($1) AND user_id=($2)",
		projectId, userId).
		Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierr.ErrNotFound
		}
		return "", fmt.Errorf("postgres query user role: %w", err)
	}

	return role, nil
}

func (r *RoleRepo) CountMembers(ctx context.Context, projectId string) (int, error) {
	var count int
	err := r.db.QueryRowContext(
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
