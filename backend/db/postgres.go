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
	if connString == "" {
		return fmt.Errorf("connect postgres: missing connection string")
	}
	pgDb, err = sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	return nil
}

func Migrate() error {
	driver, err := postgres.WithInstance(pgDb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres migrate: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("postgres migrate: %s", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("postgres migrate: %s", err)
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
		return nil, fmt.Errorf("postgres create user: %w", err)
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
			return nil, apierr.ErrResourceNotFound
		}
		return nil, fmt.Errorf("postgres find user with IDP: %w", err)
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
		return nil, fmt.Errorf("postgres create session: %w", err)
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
		return nil, fmt.Errorf("postgres create session get: %w", err)
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
		return nil, apierr.ErrResourceNotFound
	} else if err != nil {
		return nil, fmt.Errorf("postgres get active session: %w", err)
	}

	return &session, nil
}

func MakeSessionInactive(sessionId string) error {
	_, err := pgDb.Exec("UPDATE sessions "+
		"SET revoked_at = CURRENT_TIMESTAMP "+
		"WHERE id=($1)", sessionId)

	if err != nil {
		return fmt.Errorf("postgres make session inactive: %w", err)
	}

	return nil
}

func UpdateSession(sessionId string, refreshTokenEncrypted []byte, expiresAt time.Time) error {
	_, err := pgDb.Exec("UPDATE sessions "+
		"SET refresh_token_encrypted = ($1), "+
		"expires_at = ($2), "+
		"last_active_at = CURRENT_TIMESTAMP "+
		"WHERE id=($3)", refreshTokenEncrypted, expiresAt, sessionId)

	return fmt.Errorf("postgres update session: %w", err)
}

func GetUserBySub(sub string) (*models.User, error) {
	var user models.User
	err := pgDb.QueryRow("SELECT "+
		"id, idp_subject, idp_provider, username, display_name, email, avatar_url, "+
		"is_active, created_at, updated_at, last_login_at "+
		"FROM users "+
		"WHERE idp_subject=($1)",
		sub).
		Scan(&user.Id, &user.IDPSubject, &user.IDPProvider,
			&user.Username, &user.DisplayName, &user.Email, &user.AvaterURL,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginTime)
	if err != nil {
		return nil, fmt.Errorf("postgres get user by sub: %w", err)
	}

	return &user, nil
}

func CreateProject(name, description, skills, userId string) (*models.Project, error) {
	pid := uuid.NewString()
	_, err := pgDb.Exec("INSERT INTO "+
		"projects(id, name, description, skills, owner) "+
		"VALUES($1, $2, $3, $4, $5)",
		pid, name, description, skills, userId)
	if err != nil {
		return nil, fmt.Errorf("postgres create project: %w", err)
	}

	var project models.Project
	err = pgDb.QueryRow("SELECT "+
		"id, name, description, owner, skills, created_at, updated_at "+
		"FROM projects "+
		"WHERE id=($1)",
		pid).
		Scan(&project.Id, &project.Name, &project.Description,
			&project.Owner, &project.Skills, &project.CreatedAt, &project.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("postgres create project get: %w", err)
	}

	return &project, nil
}
