package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type UserRepo struct {
	db Execer
}

func NewUserRepo(db Execer) interfaces.UserRepository {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context,
	idpSubject, idpProvider, username, email string,
	displayName, avatarUrl *string) (string, error) {
	id := uuid.NewString()
	now := time.Now()

	_, err := r.db.ExecContext(ctx,
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
	err := r.db.QueryRowContext(ctx, "SELECT "+
		"id, idp_subject, idp_provider, username, display_name, email, "+
		"avatar_url, is_active, created_at, updated_at, last_login_at "+
		"FROM users "+
		"WHERE id=($1)",
		id).
		Scan(&user.Id, &user.IDPSubject, &user.IDPProvider, &user.Username,
			&user.DisplayName, &user.Email, &user.AvatarURL, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierr.ErrNotFound
		}
		return nil, fmt.Errorf("store get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetBySubject(ctx context.Context,
	idpSubject, idpProvider string) (*domain.User, error) {

	var user domain.User
	err := r.db.QueryRowContext(ctx,
		"SELECT "+
			"id, idp_subject, idp_provider, username, display_name, email, "+
			"avatar_url, is_active, created_at, updated_at, last_login_at "+
			"FROM users "+
			"WHERE idp_subject=($1) AND idp_provider=($2)",
		idpSubject, idpProvider).
		Scan(&user.Id, &user.IDPSubject, &user.IDPProvider, &user.Username,
			&user.DisplayName, &user.Email, &user.AvatarURL, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierr.ErrNotFound
		}
		return nil, fmt.Errorf("store get user with IDP: %w", err)
	}

	return &user, nil
}
