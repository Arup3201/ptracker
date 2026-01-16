package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ptracker/db"
	"github.com/ptracker/models"
	"github.com/ptracker/utils"
)

func GetExploreProjects(w http.ResponseWriter, r *http.Request) error {
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

	projectOverviews, err := db.GetExploredProjects(userId, page, limit)

	json.NewEncoder(w).Encode(HTTPSuccessResponse[models.ProjectOverviewsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &models.ProjectOverviewsResponse{
			Projects: projectOverviews,
			Page:     page,
			Limit:    limit,
			HasNext:  false,
		},
	})
	return nil
}
