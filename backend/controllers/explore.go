package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ptracker/models"
	"github.com/ptracker/services"
	"github.com/ptracker/utils"
)

type ExploreHandler struct {
	DB     *sql.DB
	UserId string
}

func (eh *ExploreHandler) GetExploreProjects(w http.ResponseWriter, r *http.Request) error {
	queryPage := r.URL.Query().Get("page")
	queryLimit := r.URL.Query().Get("limit")

	var page, limit int
	if queryPage != "" {
		var err error
		page, err = strconv.Atoi(queryPage)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Query 'page' value should be integer",
				ErrId:   ERR_INVALID_QUERY,
				Err:     fmt.Errorf("create project: validate payload: %w", err),
			}
		}
	} else {
		page = 1
	}
	if queryLimit != "" {
		var err error
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Query 'limit' value should be integer",
				ErrId:   ERR_INVALID_QUERY,
				Err:     fmt.Errorf("create project: validate payload: %w", err),
			}
		}
	} else {
		limit = 10
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get projects userId: %w", err)
	}

	exploreService := &services.ExploreService{
		DB:     eh.DB,
		UserId: userId,
	}
	projectOverviews, err := exploreService.GetExploredProjects(page, limit)
	if err != nil {
		return fmt.Errorf("service get explore projects: %w", err)
	}

	projects := []ProjectOverview{}
	for _, po := range projectOverviews {
		projects = append(projects, ProjectOverview{
			Id:          po.Id,
			Name:        po.Name,
			Description: po.Description,
			Skills:      po.Skills,
			Role:        po.Role,
			CreatedAt:   po.CreatedAt,
			UpdatedAt:   po.UpdatedAt,
		})
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectOverviewsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ProjectOverviewsResponse{
			Projects: projects,
			Page:     page,
			Limit:    limit,
			HasNext:  false,
		},
	})
	return nil
}

func (eh *ExploreHandler) GetProject(w http.ResponseWriter, r *http.Request) error {
	projectId := r.PathValue("id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Query project 'id' can't be empty",
			Err:     fmt.Errorf("get project id: empty project 'id'"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get context user: %w", err)
	}

	exploreService := &services.ExploreService{
		DB:     eh.DB,
		UserId: userId,
	}

	project, err := exploreService.GetProject(projectId)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	userStore := &models.UserStore{
		DB: eh.DB,
	}
	user, err := userStore.Get(userId)
	if err != nil {
		return fmt.Errorf("database get user: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ExploredProjectDetailsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ExploredProjectDetailsResponse{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			Skills:      project.Skills,
			Owner: Owner{
				Id:          user.Id,
				DisplayName: user.DisplayName,
				Username:    user.Username,
			},
			UnassignedTasks: project.UnassignedTasks,
			OngoingTasks:    project.OngoingTasks,
			CompletedTasks:  project.CompletedTasks,
			AbandonedTasks:  project.AbandonedTasks,
			CreatedAt:       project.CreatedAt,
			UpdatedAt:       project.UpdatedAt,
		},
	})

	return nil
}
