package models

import "time"

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
	UpdatedAt     *time.Time
	LastLoginTime time.Time
}
