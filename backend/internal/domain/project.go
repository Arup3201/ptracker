package domain

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Skills      *string   `json:"skills"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectSummary struct {
	Project
	UnassignedTasks int64 `json:"unassigned_tasks"`
	OngoingTasks    int64 `json:"ongoing_tasks"`
	CompletedTasks  int64 `json:"completed_tasks"`
	AbandonedTasks  int64 `json:"abandoned_tasks"`
}

type ProjectDetail struct {
	ProjectSummary
	Role        string `json:"role"`
	MemberCount int64  `json:"members_count"`
	Owner       Avatar `json:"owner"`
}

type ProjectPreview struct {
	Project
}

type ProjectPublicDetail struct {
	ProjectSummary
	MemberCount int64  `json:"members_count"`
	Owner       Avatar `json:"owner"`
}
