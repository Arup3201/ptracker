package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/stores"
)

type UserRepo struct {
	DB Execer
}

func NewUserRepo(db Execer) stores.UserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) Create(ctx context.Context,
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

func (r *UserRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.DB.QueryRowContext(ctx, "SELECT "+
		"id, idp_subject, idp_provider, username, display_name, email, "+
		"avatar_url, is_active, created_at, updated_at, last_login_at "+
		"FROM users "+
		"WHERE id=($1)",
		id).
		Scan(&user.Id, &user.IDPSubject, &user.IDPProvider, &user.Username,
			&user.DisplayName, &user.Email, &user.AvaterURL, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierr.ErrNotFound
		}
		return nil, fmt.Errorf("store get user: %w", err)
	}

	return &user, nil
}
