package domain

import "time"

type Project struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Skills      *string    `json:"skills"`
	Owner       string     `json:"owner"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProjectSummary struct {
	*Project
	UnassignedTasks int `json:"unassigned_tasks"`
	OngoingTasks    int `json:"ongoing_tasks"`
	CompletedTasks  int `json:"completed_tasks"`
	AbandonedTasks  int `json:"abandoned_tasks"`
}

type ProjectDetail struct {
	*ProjectSummary
	Role        string  `json:"role"`
	MemberCount int     `json:"members_count"`
	Owner       *Member `json:"owner"`
}

type PublicProjectSummary struct {
	*ProjectSummary
	MemberCount int     `json:"members_count"`
	Owner       *Member `json:"owner"`
	JoinStatus  string  `json:"join_status"`
}
