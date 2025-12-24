package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ptracker/db"
	"github.com/ptracker/utils"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Skills      string `json:"skills"`
}

func CreateProject(w http.ResponseWriter, r *http.Request) error {
	var payload CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Payload must have 'name' and 'description'",
			Err:     fmt.Errorf("create project: decode payload: %w", err),
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
		Status: "success",
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
