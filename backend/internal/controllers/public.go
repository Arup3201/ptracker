package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/utils"
)

type publicController struct {
	service interfaces.PublicService
}

func NewPublicController(service interfaces.PublicService) *publicController {
	return &publicController{
		service: service,
	}
}

func (c *publicController) ListProjects(w http.ResponseWriter, r *http.Request) error {
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

	projects, err := c.service.ListPublicProjects(r.Context(), userId)
	if err != nil {
		return fmt.Errorf("service get explore projects: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedPublicProjectsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedPublicProjectsResponse{
			Projects: projects,
			Page:     page,
			Limit:    limit,
			HasNext:  false,
		},
	})
	return nil
}

func (c *publicController) GetProject(w http.ResponseWriter, r *http.Request) error {
	projectId := r.PathValue("id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Query project 'id' can't be empty",
			Err:     fmt.Errorf("get project id: empty project 'id'"),
		}
	}

	project, err := c.service.GetPublicProject(r.Context(), projectId)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[domain.PublicProjectSummary]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   project,
	})

	return nil
}

func (c *publicController) JoinProject(w http.ResponseWriter, r *http.Request) error {

	projectId := r.PathValue("id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'id' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty 'id' provided"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get projects userId: %w", err)
	}

	err = c.service.JoinProject(r.Context(), projectId, userId)
	if err != nil {
		if errors.Is(err, apierr.ErrDuplicate) {
			return &HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Join request has already been sent",
				ErrId:   ERR_INVALID_BODY,
				Err:     fmt.Errorf("attempted for duplicate join request"),
			}
		}
		return fmt.Errorf("explore service project join request: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Join request created for the user",
	})

	return nil
}
