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

var (
	PG_HOST string
	PG_PORT string
	PG_USER string
	PG_PASS string
	PG_DB   string
)

var pgDb *sql.DB

func ConnectPostgres() error {
	var err error
	pgDb, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", PG_HOST, PG_PORT,
		PG_USER, PG_PASS, PG_DB))
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
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
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
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
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

func CreateSession(userId string, refreshTokenEncrypted []byte, userAgent, ipAddress, deviceName string,
	expireAt time.Time) (*models.Session, error) {

	sid := uuid.NewString()
	_, err := pgDb.Exec("INSERT INTO "+
		"sessions(id, user_id, refresh_token_encrypted, user_agent, "+
		"ip_address, device_name, expires_at)"+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		sid, userId, refreshTokenEncrypted, userAgent, ipAddress, deviceName, expireAt)
	if err != nil {
		return nil, err
	}

	var session models.Session
	err = pgDb.QueryRow("SELECT "+
		"id, user_id, refresh_token_encrypted, user_agent, ip_address, device_name, "+
		"created_at, last_active_at, revoked_at, expires_at "+
		"FROM sessions "+
		"WHERE id=($1)",
		sid).
		Scan(&session.Id, &session.UserId, &session.RefreshTokenEncrypted,
			&session.UserAgent, &session.IpAddress, &session.DeviceName, &session.CreatedAt, &session.LastActiveAt,
			&session.RevokedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetActiveSession(sessionId string) (*models.Session, error) {
	var session models.Session
	err := pgDb.QueryRow("SELECT "+
		"id, user_id, refresh_token_encrypted, user_agent, ip_address, device_name, "+
		"created_at, last_active_at, revoked_at, expires_at "+
		"FROM sessions "+
		"WHERE id=($1) AND revoked_at IS NULL AND expires_at>=CURRENT_TIMESTAMP",
		sessionId).
		Scan(&session.Id, &session.UserId, &session.RefreshTokenEncrypted,
			&session.UserAgent, &session.IpAddress, &session.DeviceName, &session.CreatedAt, &session.LastActiveAt,
			&session.RevokedAt, &session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, &apierr.ResourceNotFound{}
	} else if err != nil {
		return nil, err
	}

	return &session, nil
}

func MakeSessionInactive(sessionId string) error {
	_, err := pgDb.Exec("UPDATE sessions "+
		"SET revoked_at = CURRENT_TIMESTAMP "+
		"WHERE id=($1)", sessionId)

	return err
}

func UpdateSession(sessionId string, refreshTokenEncrypted []byte, expiresAt time.Time) error {
	_, err := pgDb.Exec("UPDATE sessions "+
		"SET refresh_token_encrypted = ($1), "+
		"expires_at = ($2), "+
		"last_active_at = CURRENT_TIMESTAMP "+
		"WHERE id=($3)", refreshTokenEncrypted, expiresAt, sessionId)

	return err
}
