package domain

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

type ProjectSummary struct {
	Id              string
	Name            string
	Description     *string
	Skills          *string
	Role            string
	UnassignedTasks int
	OngoingTasks    int
	CompletedTasks  int
	AbandonedTasks  int
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
