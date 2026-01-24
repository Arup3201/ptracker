package controllers

import "time"

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Skills      string `json:"skills"`
}

type Owner struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type CreatedProject struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Skills      *string    `json:"skills,omitempty"`
	Owner       Owner      `json:"owner"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
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

type ProjectSummaryResponse struct {
	ProjectSummaries []ProjectSummary `json:"projects"`
	Page             int              `json:"page"`
	Limit            int              `json:"limit"`
	HasNext          bool             `json:"has_next"`
}

type ProjectDetails struct {
	Id              string     `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Skills          *string    `json:"skills,omitempty"`
	Owner           Owner      `json:"owner"`
	Role            string     `json:"role"`
	UnassignedTasks int        `json:"unassigned_tasks"`
	OngoingTasks    int        `json:"ongoing_tasks"`
	CompletedTasks  int        `json:"completed_tasks"`
	AbandonedTasks  int        `json:"abandoned_tasks"`
	MemberCount     int        `json:"member_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

type ProjectTask struct {
	Id        string     `json:"id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ProjectTasksResponse struct {
	ProjectTasks []ProjectTask `json:"tasks"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
	HasNext      bool          `json:"has_next"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status" validate:"required"`
}

type CreatedProjectTask struct {
	Id          string     `json:"id"`
	ProjectId   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProjectTaskDetails struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProjectOverview struct {
	Id          string     `json:"id"`
	Name        string     `json:"title"`
	Description *string    `json:"description"`
	Skills      *string    `json:"skills"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProjectOverviewsResponse struct {
	Projects []ProjectOverview `json:"projects"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	HasNext  bool              `json:"has_next"`
}

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
}

type JoinRequest struct {
	ProjectId string `json:"project_id"`
	User      User   `json:"user"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type JoinRequestsResponse struct {
	Requests []JoinRequest `json:"join_requests"`
}

type ExploredProjectDetailsResponse struct {
	Id              string     `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	Skills          *string    `json:"skills"`
	Owner           Owner      `json:"owner"`
	JoinStatus      string     `json:"join_status"`
	UnassignedTasks int        `json:"unassigned_tasks"`
	OngoingTasks    int        `json:"ongoing_tasks"`
	CompletedTasks  int        `json:"completed_tasks"`
	AbandonedTasks  int        `json:"abandoned_tasks"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
