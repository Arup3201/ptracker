package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type SessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) interfaces.SessionRepository {
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
	session := models.Session{
		ID:                    id,
		UserID:                userId,
		RefreshTokenEncrypted: encryptedToken,
		UserAgent:             userAgent,
		IpAddress:             ipAddress,
		DeviceName:            deviceName,
		ExpiresAt:             expireAt,
	}
	err := gorm.G[models.Session](r.db).Create(ctx, &session)
	if err != nil {
		return "", fmt.Errorf("gorm create: %w", err)
	}

	return session.ID, nil
}

func (r *SessionRepo) Get(ctx context.Context, id string) (domain.Session, error) {

	session, err := gorm.G[models.Session](r.db).Where("id = ? AND revoked_at IS NULL",
		id).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return session.ToSessionDomain(), apierr.ErrNotFound
	} else if err != nil {
		return session.ToSessionDomain(), fmt.Errorf("gorm query: %w", err)
	}

	return session.ToSessionDomain(), nil
}

func (r *SessionRepo) Revoke(ctx context.Context, id string) error {

	session, err := r.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("session get: %w", err)
	}

	now := time.Now().UTC()
	session.RevokedAt = &now

	err = r.db.Save(&session).Error
	if err != nil {
		return fmt.Errorf("gorm db save: %w", err)
	}

	return nil
}

func (r *SessionRepo) Update(ctx context.Context, id string,
	refreshTokenEncrypted []byte,
	expiresAt time.Time) error {

	session, err := r.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("session get: %w", err)
	}

	session.RefreshTokenEncrypted = refreshTokenEncrypted
	session.ExpiresAt = expiresAt

	err = r.db.Save(&session).Error
	if err != nil {
		return fmt.Errorf("gorm db save: %w", err)
	}

	return nil
}
