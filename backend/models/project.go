package models

import "time"

type Project struct {
	Id          string
	Name        string
	Description *string
	Skills      *string
	Owner       string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
