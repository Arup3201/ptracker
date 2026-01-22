package controllers

import (
	"time"

	"github.com/ptracker/models"
)

type UserStore interface {
	Create(idpSubject, idpProvider, username,
		displayName, email, avatarUrl string) (string, error)
	GetBySubject(idpSubject, idpProvider string) (models.User, error)
}

type SessionStore interface {
	Create(userId string,
		refreshTokenEncrypted []byte,
		userAgent, ipAddress, deviceName string,
		expireAt time.Time) (string, error)
	Get(id string) (models.Session, error)
	Update(id string,
		tokenEncrypted []byte,
		expiresAt time.Time) error
	Revoke(id string) error
}

type ProjectStore interface {
	Create(name, description, skills string) (string, error)
	Get(id string) (models.Project, error)
	All(page, limit int) ([]models.Project, error)
	Count() (int, error)
}

type TaskStore interface {
	Create(title, description, status string) (string, error)
	Get(id string) (models.ProjectTask, error)
	All() ([]models.ProjectTask, error)
	Count() (int, error)
}
