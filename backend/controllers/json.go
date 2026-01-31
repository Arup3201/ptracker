package controllers

import (
	"github.com/ptracker/domain"
)

type CreateProjectRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Skills      *string `json:"skills"`
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

type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
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

type UpdateJoinRequest struct {
	UserId     string `json:"user_id" validate:"required"`
	JoinStatus string `json:"join_status" validate:"required"`
}

type MembersResponse struct {
	Members []*domain.Member `json:"members"`
}
