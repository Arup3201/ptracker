package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
