package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/apierr"
)

type User struct {
	Id            string
	IDPSubject    string
	IDPProvider   string
	Username      string
	DisplayName   string
	Email         string
	AvaterURL     string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     *time.Time // nullable
	LastLoginTime time.Time
}

type UserStore struct {
	DB *sql.DB
}

func (us *UserStore) Create(idpSubject, idpProvider, username,
	displayName, email, avatarUrl string) (string, error) {
	uId := uuid.NewString()
	_, err := us.DB.Exec("INSERT INTO users"+
		"(id, idp_subject, idp_provider, username, display_name, email, "+
		"avatar_url) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7)",
		uId, idpSubject, idpProvider, username, displayName, email, avatarUrl)
	if err != nil {
		return "", fmt.Errorf("store create user: %w", err)
	}

	return uId, nil
}

func (us *UserStore) Get(id string) (User, error) {
	var user User
	err := us.DB.QueryRow("SELECT "+
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
			return user, apierr.ErrResourceNotFound
		}
		return user, fmt.Errorf("store get user: %w", err)
	}

	return user, nil
}

func (us *UserStore) GetBySubject(idpSubject, idpProvider string) (User, error) {
	var user User
	err := us.DB.QueryRow("SELECT "+
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
			return user, apierr.ErrResourceNotFound
		}
		return user, fmt.Errorf("store get user with IDP: %w", err)
	}

	return user, nil
}
