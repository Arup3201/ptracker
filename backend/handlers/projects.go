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

	data, err := db.GetAllProjects(userId, page, limit)
	if err != nil {
		return fmt.Errorf("get projects from db: %w", err)
	}

	projectSummaries := []models.ProjectSummary{}
	for _, row := range data {
		projectSummaries = append(projectSummaries, models.ProjectSummary{
			Id:   row.Id,
			Name: row.Name,
		})
	}

	cnt, err := db.GetProjectsCount(userId)
	if err != apierr.ErrResourceNotFound && err != nil {
		return fmt.Errorf("get projects from db: %w", err)
	}

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: HTTPData{
			"projects": projectSummaries,
			"page":     page,
			"limit":    limit,
			"total":    cnt,
			"has_next": hasNext,
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
			Err:     fmt.Errorf("create project: decode payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' is missing",
			Err:     fmt.Errorf("create project: validate payload: %w", err),
		}
	}

	if payload.Name == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' can't be empty",
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

	json.NewEncoder(w).Encode(HTTPSuccessResponse{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: HTTPData{
			"id":          project.Id,
			"name":        project.Name,
			"description": project.Description,
			"skills":      project.Skills,
			"owner":       project.Owner,
			"created_at":  project.CreatedAt,
			"updated_at":  project.UpdateAt,
		},
	})
	return nil
}
