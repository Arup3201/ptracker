package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type userRepo struct {
	DB Execer
}

func NewUserRepo(db Execer) *userRepo {
	return &userRepo{
		DB: db,
	}
}

func (r *userRepo) Create(ctx context.Context,
	idpSubject, idpProvider, username, email string,
	displayName, avatarUrl *string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.DB.ExecContext(ctx,
		"INSERT INTO users"+
			"(id, idp_subject, idp_provider, username, display_name, email, "+
			"avatar_url, created_at, updated_at) "+
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $8)",
		id, idpSubject, idpProvider, username, displayName, email, avatarUrl, now)
	if err != nil {
		return "", fmt.Errorf("store create user: %w", err)
	}

	return id, nil
}
