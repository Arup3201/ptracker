package controllers

import "github.com/ptracker/internal/domain"

type CreateProjectRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Skills      *string `json:"skills"`
}

type CreateTaskRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Assignees   []string `json:"assignees" validate:"required"`
	Status      string   `json:"status" validate:"required"`
}

type CreateTaskResponse struct {
	Id       string   `json:"task_id"`
	Warnings []string `json:"warnings"`
}

type UpdateJoinRequest struct {
	UserId     string `json:"user_id" validate:"required"`
	JoinStatus string `json:"join_status" validate:"required"`
}

type UpdateTaskRequest struct {
	Title             string   `json:"title" validate:"required"`
	Description       string   `json:"description" validate:"required"`
	Status            string   `json:"status" validate:"required"`
	AssigneesToAdd    []string `json:"assignees_to_add" validate:"required"`
	AssigneesToRemove []string `json:"assignees_to_remove" validate:"required"`
}

type ListedPrivateProjectsResponse struct {
	ProjectSummaries []*domain.PrivateProjectListed `json:"projects"`
	Page             int                            `json:"page"`
	Limit            int                            `json:"limit"`
	HasNext          bool                           `json:"has_next"`
}

type ListedTasksResponse struct {
	ProjectTasks []*domain.TaskListed `json:"tasks"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	HasNext      bool                 `json:"has_next"`
}

type ListedPublicProjectsResponse struct {
	Projects []*domain.PublicProjectListed `json:"projects"`
	Page     int                           `json:"page"`
	Limit    int                           `json:"limit"`
	HasNext  bool                          `json:"has_next"`
}

type ListedJoinRequestsResponse struct {
	Requests []*domain.JoinRequestListed `json:"join_requests"`
}

type ListedMembersResponse struct {
	Members []*domain.Member `json:"members"`
}
