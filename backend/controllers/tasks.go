package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/models"
	"github.com/ptracker/utils"
)

type TaskHandler struct {
	DB *sql.DB
}

func (th *TaskHandler) All(w http.ResponseWriter, r *http.Request) error {
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

	roleStore := &models.RoleStore{
		DB: th.DB,
	}
	taskStore := &models.TaskStore{
		DB:        th.DB,
		ProjectId: projectId,
		UserId:    userId,
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

	tasks, err := taskStore.All()
	if err != nil {
		return fmt.Errorf("database get project tasks: %w", err)
	}

	responseTasks := []ProjectTask{}
	for _, t := range tasks {
		responseTasks = append(responseTasks, ProjectTask{
			Id:        t.Id,
			Title:     t.Title,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	cnt, err := taskStore.Count()
	if err != nil {
		return fmt.Errorf("database count project task: %w", err)
	}

	hasNext := cnt > page*limit

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectTasksResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ProjectTasksResponse{
			ProjectTasks: responseTasks,
			Page:         page,
			Limit:        limit,
			HasNext:      hasNext,
		},
	})

	return nil
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) error {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Query `project_id' can't be empty",
			Err:     fmt.Errorf("empty 'project_id'"),
		}
	}

	var payload CreateTaskRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task accepts title, description, assignee and status only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("decode task payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task 'title' or 'status' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("task payload validation: %w", err),
		}
	}

	if payload.Title == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task 'title' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty task 'title' provided"),
		}
	}
	if payload.Status == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task 'status' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty task 'status' provided"),
		}
	}
	if !slices.Contains(models.TASK_STATUS, payload.Status) {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task 'status' is invalid, example: " + strings.Join(models.TASK_STATUS, ", "),
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("task 'status' invalid value provided"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	roleStore := &models.RoleStore{
		DB: th.DB,
	}
	taskStore := &models.TaskStore{
		DB:        th.DB,
		ProjectId: projectId,
		UserId:    userId,
	}

	access, err := roleStore.CanEdit()
	if err != nil {
		return fmt.Errorf("database check write permission: %w", err)
	}
	if !access {
		return &HTTPError{
			Code:    http.StatusForbidden,
			ErrId:   ERR_ACCESS_DENIED,
			Message: "User is not the owner of this project",
		}
	}

	taskId, err := taskStore.Create(payload.Title, payload.Description, payload.Status)
	if err != nil {
		return fmt.Errorf("db create task: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &taskId,
	})

	return nil
}

func (th *TaskHandler) Get(w http.ResponseWriter, r *http.Request) error {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Query `project_id' can't be empty",
			Err:     fmt.Errorf("empty 'project_id'"),
		}
	}

	taskId := r.PathValue("task_id")
	if taskId == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Query `task_id' can't be empty",
			Err:     fmt.Errorf("empty 'task_id'"),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	roleStore := &models.RoleStore{
		DB: th.DB,
	}
	taskStore := &models.TaskStore{
		DB:        th.DB,
		ProjectId: projectId,
		UserId:    userId,
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

	task, err := taskStore.Get(taskId)
	if err != nil {
		return fmt.Errorf("db get task: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectTaskDetails]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ProjectTaskDetails{
			Id:          task.Id,
			Title:       task.Title,
			Description: &task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	})

	return nil
}
