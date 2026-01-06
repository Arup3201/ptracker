package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/apierr"
	"github.com/ptracker/db"
	"github.com/ptracker/models"
	"github.com/ptracker/utils"
)

func GetAllProjects(w http.ResponseWriter, r *http.Request) error {
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

	projectSummaries, err := db.GetAllProjects(userId, page, limit)
	if err != nil {
		return fmt.Errorf("get projects from db: %w", err)
	}

	cnt, err := db.GetProjectsCount(userId)
	if err != apierr.ErrResourceNotFound && err != nil {
		return fmt.Errorf("get projects from db: %w", err)
	}

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse[models.ProjectSummaryResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &models.ProjectSummaryResponse{
			ProjectSummaries: projectSummaries,
			Page:             page,
			Limit:            limit,
			HasNext:          hasNext,
		},
	})

	return nil
}

func CreateProject(w http.ResponseWriter, r *http.Request) error {
	var payload models.CreateProjectRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project must have 'name' with optional 'description' and 'skills' fields only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("create project: decode payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("create project: validate payload: %w", err),
		}
	}

	if payload.Name == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("create project: empty project 'name' provided"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("create project: get userId: %w", err)
	}

	project, err := db.CreateProject(payload.Name, payload.Description, payload.Skills, userId)
	if err != nil {
		return fmt.Errorf("create project: create project: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[models.CreatedProject]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   project,
	})
	return nil
}

func GetProject(w http.ResponseWriter, r *http.Request) error {
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

	project, err := db.GetProjectDetails(userId, projectId)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[models.ProjectDetails]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   project,
	})

	return nil
}
