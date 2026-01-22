package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/apierr"
)

type Session struct {
	Id                    string
	UserId                string
	RefreshTokenEncrypted []byte
	UserAgent             string
	IpAddress             string
	DeviceName            string
	CreatedAt             time.Time
	LastActiveAt          time.Time
	RevokedAt             *time.Time // nullable
	ExpiresAt             time.Time
}

type SessionStore struct {
	DB *sql.DB
}

func (ss *SessionStore) Create(userId string, refreshTokenEncrypted []byte, userAgent, ipAddress, deviceName string,
	expireAt time.Time) (string, error) {

	sId := uuid.NewString()
	_, err := ss.DB.Exec("INSERT INTO "+
		"sessions(id, user_id, refresh_token_encrypted, user_agent, "+
		"ip_address, device_name, expires_at)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		sId, userId, refreshTokenEncrypted, userAgent, ipAddress, deviceName, expireAt)
	if err != nil {
		return "", fmt.Errorf("postgres create session: %w", err)
	}

	var session Session
	err = ss.DB.QueryRow("SELECT "+
		"id, user_id, refresh_token_encrypted, user_agent, ip_address, device_name, "+
		"created_at, last_active_at, revoked_at, expires_at "+
		"FROM sessions "+
		"WHERE id=($1)",
		sId).
		Scan(&session.Id, &session.UserId, &session.RefreshTokenEncrypted,
			&session.UserAgent, &session.IpAddress, &session.DeviceName, &session.CreatedAt, &session.LastActiveAt,
			&session.RevokedAt, &session.ExpiresAt)
	if err != nil {
		return "", fmt.Errorf("postgres create session get: %w", err)
	}

	return sId, nil
}

func (ss *SessionStore) Get(sessionId string) (Session, error) {
	var session Session
	err := ss.DB.QueryRow("SELECT "+
		"id, user_id, refresh_token_encrypted, user_agent, ip_address, device_name, "+
		"created_at, last_active_at, revoked_at, expires_at "+
		"FROM sessions "+
		"WHERE id=($1) AND revoked_at IS NULL AND expires_at>=CURRENT_TIMESTAMP",
		sessionId).
		Scan(&session.Id, &session.UserId, &session.RefreshTokenEncrypted,
			&session.UserAgent, &session.IpAddress, &session.DeviceName, &session.CreatedAt, &session.LastActiveAt,
			&session.RevokedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return session, apierr.ErrResourceNotFound
	} else if err != nil {
		return session, fmt.Errorf("postgres get active session: %w", err)
	}

	return session, nil
}

func (ss *SessionStore) Revoke(sessionId string) error {
	_, err := ss.DB.Exec("UPDATE sessions "+
		"SET revoked_at = CURRENT_TIMESTAMP "+
		"WHERE id=($1)", sessionId)

	if err != nil {
		return fmt.Errorf("postgres make session inactive: %w", err)
	}

	return nil
}

func (ss *SessionStore) Update(sessionId string,
	refreshTokenEncrypted []byte,
	expiresAt time.Time) error {
	_, err := ss.DB.Exec("UPDATE sessions "+
		"SET refresh_token_encrypted = ($1), "+
		"expires_at = ($2), "+
		"last_active_at = CURRENT_TIMESTAMP "+
		"WHERE id=($3)", refreshTokenEncrypted, expiresAt, sessionId)

	if err != nil {
		return fmt.Errorf("postgres update session: %w", err)
	}
	return nil
}
