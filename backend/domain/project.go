package domain

import "time"

type PrivateProject struct {
	Id          string
	Name        string
	Description *string
	Skills      *string
	Owner       string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type ProjectSummary struct {
	PrivateProject
	UnassignedTasks int
	OngoingTasks    int
	CompletedTasks  int
	AbandonedTasks  int
}

type ProjectDetail struct {
	ProjectSummary
	Role        string
	MemberCount int
	Owner       *Member
}
