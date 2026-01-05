package models

import "time"

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Skills      string `json:"skills"`
}

type ProjectSummary struct {
	Id              string     `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	Skills          *string    `json:"skills"`
	Role            string     `json:"role"`
	UnassignedTasks int        `json:"unassigned_tasks"`
	OngoingTasks    int        `json:"ongoing_tasks"`
	CompletedTasks  int        `json:"completed_tasks"`
	AbandonedTasks  int        `json:"abandoned_tasks"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type Owner struct {
	Id          string
	Username    string
	DisplayName string
}

type ProjectDetails struct {
	Id              string
	Name            string
	Description     *string
	Skills          *string
	Owner           Owner
	Role            string
	UnassignedTasks int
	OngoingTasks    int
	CompletedTasks  int
	AbandonedTasks  int
	MemberCount     int
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
