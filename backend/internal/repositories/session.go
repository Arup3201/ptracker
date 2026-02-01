package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
)

type SessionRepo struct {
	db interfaces.Execer
}

func NewSessionRepo(db interfaces.Execer) interfaces.SessionRepository {
	return &SessionRepo{
		db: db,
	}
}

func (r *SessionRepo) Create(ctx context.Context,
	userId string,
	encryptedToken []byte,
	userAgent, ipAddress, deviceName string,
	expireAt time.Time) (string, error) {

	id := uuid.NewString()
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO "+
			"sessions(id, user_id, refresh_token_encrypted, user_agent, "+
			"ip_address, device_name, expires_at)"+
			"VALUES($1, $2, $3, $4, $5, $6, $7)",
		id, userId, encryptedToken, userAgent, ipAddress, deviceName, expireAt)
	if err != nil {
		return "", fmt.Errorf("postgres create session: %w", err)
	}

	return id, nil
}

func (r *SessionRepo) Get(ctx context.Context, id string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.QueryRowContext(ctx,
		"SELECT "+
			"id, user_id, refresh_token_encrypted, "+
			"user_agent, ip_address, device_name, "+
			"created_at, last_active_at, revoked_at, expires_at "+
			"FROM sessions "+
			"WHERE id=($1)",
		id).
		Scan(&session.Id, &session.UserId, &session.RefreshTokenEncrypted,
			&session.UserAgent, &session.IpAddress, &session.DeviceName,
			&session.CreatedAt, &session.LastActive,
			&session.RevokedAt, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("postgres create session get: %w", err)
	}

	return &session, nil
}

func (r *SessionRepo) Revoke(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sessions "+
			"SET revoked_at = CURRENT_TIMESTAMP "+
			"WHERE id=($1)", id)

	if err != nil {
		return fmt.Errorf("postgres make session inactive: %w", err)
	}

	return nil
}

func (r *SessionRepo) Update(ctx context.Context, sessionId string,
	refreshTokenEncrypted []byte,
	expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE sessions "+
			"SET refresh_token_encrypted = ($1), "+
			"expires_at = ($2), "+
			"last_active_at = CURRENT_TIMESTAMP "+
			"WHERE id=($3)", refreshTokenEncrypted, expiresAt, sessionId)

	if err != nil {
		return fmt.Errorf("postgres update session: %w", err)
	}
	return nil
}
