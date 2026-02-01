package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/utils"
)

type taskController struct {
	service interfaces.TaskService
}

func NewTaskController(service interfaces.TaskService) *taskController {
	return &taskController{
		service: service,
	}
}

func (c *taskController) ListTasks(w http.ResponseWriter, r *http.Request) error {

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

	tasks, err := c.service.ListTasks(r.Context(), projectId, userId)
	if err != nil {
		return fmt.Errorf("service list tasks: %w", err)
	}

	// TODO: count total tasks and check has next
	hasNext := false

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedTasksResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedTasksResponse{
			ProjectTasks: tasks,
			Page:         page,
			Limit:        limit,
			HasNext:      hasNext,
		},
	})

	return nil
}

func (c *taskController) CreateTask(w http.ResponseWriter, r *http.Request) error {

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

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	taskId, err := c.service.CreateTask(
		r.Context(),
		projectId,
		payload.Title,
		payload.Description,
		userId)
	if err != nil {
		return fmt.Errorf("service create task: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &taskId,
	})

	return nil
}

func (c *taskController) GetTask(w http.ResponseWriter, r *http.Request) error {

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

	task, err := c.service.GetTask(
		r.Context(),
		projectId,
		taskId,
		userId)
	if err != nil {
		return fmt.Errorf("service get task: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[domain.Task]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   task,
	})

	return nil
}
