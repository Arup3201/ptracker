package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/apierr"
	"github.com/ptracker/models"
	"github.com/ptracker/utils"
)

type ProjectHandler struct {
	DB *sql.DB
}

func (ph *ProjectHandler) All(w http.ResponseWriter, r *http.Request) error {
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

	projectStore := &models.ProjectStore{
		DB:     ph.DB,
		UserId: userId,
	}
	projectSummaries, err := projectStore.All(page, limit)
	if err != nil {
		return fmt.Errorf("get projects from store: %w", err)
	}

	summaries := []ProjectSummary{}
	for _, ps := range projectSummaries {
		summaries = append(summaries, ProjectSummary{
			Id:              ps.Id,
			Name:            ps.Name,
			Description:     ps.Description,
			Skills:          ps.Skills,
			Role:            ps.Role,
			UnassignedTasks: ps.UnassignedTasks,
			OngoingTasks:    ps.OngoingTasks,
			CompletedTasks:  ps.CompletedTasks,
			AbandonedTasks:  ps.AbandonedTasks,
			CreatedAt:       ps.CreatedAt,
			UpdatedAt:       ps.UpdatedAt,
		})
	}

	cnt, err := projectStore.Count()
	if err != apierr.ErrResourceNotFound && err != nil {
		return fmt.Errorf("get projects from store: %w", err)
	}

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectSummaryResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ProjectSummaryResponse{
			ProjectSummaries: summaries,
			Page:             page,
			Limit:            limit,
			HasNext:          hasNext,
		},
	})

	return nil
}

func (ph *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var payload CreateProjectRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project must have 'name' with optional 'description' and 'skills' fields only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("store create project: decode payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("store create project: validate payload: %w", err),
		}
	}

	if payload.Name == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty project 'name' provided"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	projectStore := &models.ProjectStore{
		DB:     ph.DB,
		UserId: userId,
	}
	projectId, err := projectStore.Create(payload.Name, payload.Description, payload.Skills)
	if err != nil {
		return fmt.Errorf("store create project: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &projectId,
	})
	return nil
}

func (ph *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) error {
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

	projectStore := &models.ProjectStore{
		DB:     ph.DB,
		UserId: userId,
	}
	roleStore := &models.RoleStore{
		DB:        ph.DB,
		ProjectId: projectId,
		UserId:    userId,
	}
	userStore := &models.UserStore{
		DB: ph.DB,
	}

	access, err := roleStore.CanAccess()
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

	project, err := projectStore.Get(projectId)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	memberCount, err := projectStore.CountMembers(projectId)
	if err != nil {
		return fmt.Errorf("database get project member count: %w", err)
	}

	role, err := roleStore.Get()
	if err != nil {
		return fmt.Errorf("database get user role in project: %w", err)
	}

	user, err := userStore.Get(userId)
	if err != nil {
		return fmt.Errorf("database get user: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectDetails]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ProjectDetails{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			Skills:      project.Skills,
			Owner: Owner{
				Id:          user.Id,
				DisplayName: user.DisplayName,
				Username:    user.Username,
			},
			Role:            role.Role,
			MemberCount:     memberCount,
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
