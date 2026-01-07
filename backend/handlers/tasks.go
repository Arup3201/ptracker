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

func GetProjectTasks(w http.ResponseWriter, r *http.Request) error {
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

	projectId := r.PathValue("project_id")
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

	access, err := db.CanAccess(userId, projectId)
	if err != nil {
		return fmt.Errorf("database check access: %w", err)
	}
	if !access {
		return &HTTPError{
			Code:    http.StatusForbidden,
			ErrId:   ERR_ACCESS_DENIED,
			Message: "User is not a member of the project",
		}
	}

	tasks, err := db.GetProjectTasks(projectId)
	if err != nil {
		return fmt.Errorf("database get project tasks: %w", err)
	}

	cnt, err := db.GetProjectTaskCount(projectId)
	if err != nil {
		return fmt.Errorf("database count project task: %w", err)
	}

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse[models.ProjectTasksResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &models.ProjectTasksResponse{
			ProjectTasks: tasks,
			Page:         page,
			Limit:        limit,
			HasNext:      hasNext,
		},
	})

	return nil
}
