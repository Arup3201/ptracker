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
	*Project
	UnassignedTasks int
	OngoingTasks    int
	CompletedTasks  int
	AbandonedTasks  int
}

type ProjectDetail struct {
	*ProjectSummary
	Role        string
	MemberCount int
	Owner       *Member
}

type PublicProjectSummary struct {
	*ProjectSummary
	MemberCount int
	Owner       *Member
	JoinStatus  string
}
