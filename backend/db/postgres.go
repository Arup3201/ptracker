package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/google/uuid"
	"github.com/ptracker/apierr"
	"github.com/ptracker/models"
)

var pgDb *sql.DB

func ConnectPostgres(connString string) error {
	var err error
	pgDb, err = sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("[ERROR] postgres Open: %s", err)
	}

	return nil
}

func Migrate() error {
	driver, err := postgres.WithInstance(pgDb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("[ERROR] postgres WithInstance: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("[ERROR] postgres NewWithDatabaseInstance: %s", err)
	}
	if err = m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		fmt.Printf("[ERROR] postgres migration error: %s", err)
	}

	return nil
}

func CreateUser(idpSubject, idpProvider, username,
	displayName, email, avatarUrl string) (*models.User, error) {
	uId := uuid.NewString()
	_, err := pgDb.Exec("INSERT INTO users"+
		"(id, idp_subject, idp_provider, username, display_name, email, "+
		"avatar_url) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
		uId, idpSubject, idpProvider, username, displayName, email, avatarUrl)
	if err != nil {
		return nil, err
	}

	return FindUserWithIdp(idpSubject, idpProvider)
}

func FindUserWithIdp(idpSubject, idpProvider string) (*models.User, error) {
	var user models.User
	err := pgDb.QueryRow("SELECT "+
		"id, idp_subject, idp_provider, username, display_name, email, "+
		"avatar_url, is_active, created_at, updated_at, last_login_at "+
		"FROM users "+
		"WHERE idp_subject=($1) AND idp_provider=($2)",
		idpSubject, idpProvider).
		Scan(&user.Id, &user.IDPSubject, &user.IDPProvider, &user.Username,
			&user.DisplayName, &user.Email, &user.AvaterURL, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apierr.ResourceNotFound{}
		}
		return nil, err
	}

	return &user, nil
}

func CreateSession(userId, refreshTokenHash, userAgent, ipAddress, deviceName string,
	expireAt time.Time) (string, error) {

	sid := uuid.NewString()
	_, err := pgDb.Exec("INSERT INTO "+
		"sessions(id, user_id, refresh_token_hash, user_agent, "+
		"ip_address, device_name, expires_at)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		sid, userId, refreshTokenHash, userAgent, ipAddress, deviceName, expireAt)
	if err != nil {
		return "", err
	}

	return sid, nil
}

func GetSession(sessionId string) (*models.Session, error) {
	var session models.Session
	err := pgDb.QueryRow("SELECT "+
		"id, user_id, refresh_token_hash, user_agent, ip_address, device_name, "+
		"created_at, last_active_at, revoked_at, expires_at "+
		"FROM sessions"+
		"WHERE id=($1)", sessionId).Scan(&session.Id, &session.UserId, &session.RefreshTokenHash,
		&session.UserAgent, &session.IpAddress, &session.DeviceName, &session.CreatedAt, &session.LastActiveAt,
		&session.RevokedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, &apierr.ResourceNotFound{}
	}

	return &session, nil
}
