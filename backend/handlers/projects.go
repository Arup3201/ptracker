package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/db"
	"github.com/ptracker/utils"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Skills      string `json:"skills"`
}

func CreateProject(w http.ResponseWriter, r *http.Request) error {
	var payload CreateProjectRequest

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
