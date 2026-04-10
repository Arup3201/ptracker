package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type commentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) interfaces.CommentRepository {
	return &commentRepo{
		db: db,
	}
}

func (r *commentRepo) Create(ctx context.Context,
	projectId, taskId, userId string,
	content string) (string, error) {

	var err error

	id := uuid.NewString()
	comment := models.Comment{
		ID:        id,
		ProjectID: projectId,
		TaskID:    taskId,
		UserID:    userId,
		Content:   content,
	}

	err = gorm.G[models.Comment](r.db).Create(ctx, &comment)
	if err != nil {
		return "", fmt.Errorf("gorm create: %w", err)
	}

	return comment.ID, nil
}
