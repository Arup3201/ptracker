package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ptracker/internal/interfaces"
)

type commentRepo struct {
	db interfaces.Execer
}

func NewCommentRepo(db interfaces.Execer) interfaces.CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) Create(ctx context.Context,
	projectId, taskId, userId string,
	comment string) (string, error) {

	id := uuid.NewString()

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO 
		comments (id, project_id, task_id, user_id, 
		content, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		id, projectId, taskId, userId, comment)
	if err != nil {
		return "", err
	}
	return id, nil
}
