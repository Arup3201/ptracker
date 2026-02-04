package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/utils"
)

type projectController struct {
	service interfaces.ProjectService
}

func NewProjectController(service interfaces.ProjectService) interfaces.ProjectController {
	return &projectController{
		service: service,
	}
}

func (c *projectController) List(w http.ResponseWriter, r *http.Request) error {
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

	summaries, err := c.service.ListProjects(r.Context(), userId)
	if err != nil {
		return fmt.Errorf("get projects from store: %w", err)
	}

	// TODO: count projects where the user is part of
	cnt := 0

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedPrivateProjectsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedPrivateProjectsResponse{
			ProjectSummaries: summaries,
			Page:             page,
			Limit:            limit,
			HasNext:          hasNext,
		},
	})

	return nil
}

func (c *projectController) Create(w http.ResponseWriter, r *http.Request) error {

	var payload CreateProjectRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project must have 'name' with optional 'description' and 'skills' fields only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("decode payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("validate payload: %w", err),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	projectId, err := c.service.CreateProject(r.Context(),
		payload.Name, payload.Description, payload.Skills, userId)
	if err != nil {
		return fmt.Errorf("store create project: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &projectId,
	})
	return nil
}

func (c *projectController) Get(w http.ResponseWriter, r *http.Request) error {

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

	project, err := c.service.GetPrivateProject(r.Context(), projectId, userId)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[domain.ProjectDetail]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   project,
	})

	return nil
}

func (c *projectController) ListJoinRequests(w http.ResponseWriter, r *http.Request) error {

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

	joinRequests, err := c.service.ListJoinRequests(r.Context(), projectId, userId)
	if err != nil {
		return fmt.Errorf("explore service get join request: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedJoinRequestsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedJoinRequestsResponse{
			Requests: joinRequests,
		},
	})

	return nil
}

func (c *projectController) RespondToJoinRequests(w http.ResponseWriter, r *http.Request) error {

	projectId := r.PathValue("id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'id' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty 'id' provided"),
		}
	}

	var payload UpdateJoinRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Operation needs two fields: user_id and join_status",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("join request update payload decode: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Either user_id or join_status is missing or has invalid type",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("join request update payload validate: %w", err),
		}
	}

	if payload.UserId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Payload 'user_id' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty 'user_id' provided"),
		}
	}
	if payload.JoinStatus == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Payload 'join_status' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty 'join_status' provided"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get projects userId: %w", err)
	}

	err = c.service.RespondToJoinRequests(
		r.Context(),
		projectId,
		userId,
		payload.UserId,
		payload.JoinStatus,
	)
	if err != nil {
		switch err {
		case apierr.ErrInvalidValue:
			return &HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Payload 'join_status' has invalid value",
				ErrId:   ERR_INVALID_BODY,
				Err:     fmt.Errorf("invalid 'join_status' value provided"),
			}
		default:
			return fmt.Errorf("service update join request status: %w", err)
		}
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Join request status updated",
	})

	return nil
}

func (c *projectController) ListMembers(w http.ResponseWriter, r *http.Request) error {

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

	members, err := c.service.GetProjectMembers(r.Context(),
		projectId,
		userId)
	if err != nil {
		return fmt.Errorf("service members: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedMembersResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedMembersResponse{
			Members: members,
		},
	})

	return nil
}

func (c *projectController) ListRecentlyCreatedProjects(w http.ResponseWriter, r *http.Request) error {

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get projects userId: %w", err)
	}

	projects, err := c.service.ListRecentlyCreatedProjects(r.Context(),
		userId)
	if err != nil {
		return fmt.Errorf("service list recently created projects: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedRecentProjectsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedRecentProjectsResponse{
			Projects: projects,
		},
	})

	return nil
}

func (c *projectController) ListRecentlyJoinedProjects(w http.ResponseWriter, r *http.Request) error {
	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get projects userId: %w", err)
	}

	projects, err := c.service.ListRecentlyJoinedProjects(r.Context(),
		userId)
	if err != nil {
		return fmt.Errorf("service list recently created projects: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedRecentProjectsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedRecentProjectsResponse{
			Projects: projects,
		},
	})

	return nil
}
