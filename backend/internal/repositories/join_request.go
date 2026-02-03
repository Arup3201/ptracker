package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/interfaces"
)

type JoinRequestRepo struct {
	db interfaces.Execer
}

func NewJoinRequestRepo(db interfaces.Execer) interfaces.JoinRequestRepository {
	return &JoinRequestRepo{
		db: db,
	}
}

func (r *JoinRequestRepo) Create(ctx context.Context,
	projectId, userId, joinStatus string) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx,
		"INSERT INTO join_requests(user_id, project_id, status, created_at, updated_at) "+
			"VALUES($1, $2, $3, $4, $4)", userId, projectId, joinStatus, now)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return apierr.ErrDuplicate
		}
		return fmt.Errorf("postgres join project: %w", err)
	}

	return nil
}

func (r *JoinRequestRepo) Get(ctx context.Context, projectId, userId string) (string, error) {
	var status string
	err := r.db.QueryRowContext(
		ctx,
		"SELECT "+
			"status "+
			"FROM join_requests "+
			"WHERE project_id=($1) AND user_id=($2)",
		projectId, userId,
	).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierr.ErrNotFound
		}
		return "", fmt.Errorf("query row context: %w", err)
	}

	return status, nil
}

func (r *JoinRequestRepo) Update(ctx context.Context, projectId, userId, joinStatus string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE join_requests "+
		"SET status=($1) "+
		"WHERE project_id=($2) AND user_id=($3)", joinStatus, projectId, userId)

	if err != nil {
		if strings.Contains(err.Error(), "invalid input value") {
			return apierr.ErrInvalidValue
		}
		return fmt.Errorf("service update join request query: %w", err)
	}

	return nil
}
