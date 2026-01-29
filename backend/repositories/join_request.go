package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ptracker/apierr"
	"github.com/ptracker/interfaces"
)

type JoinRequestRepo struct {
	db Execer
}

func NewJoinRequestRepo(db Execer) interfaces.JoinRequestRepository {
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
