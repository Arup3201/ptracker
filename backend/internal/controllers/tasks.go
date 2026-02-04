package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/utils"
)

type taskController struct {
	service interfaces.TaskService
}

func NewTaskController(service interfaces.TaskService) interfaces.TaskController {
	return &taskController{
		service: service,
	}
}

func (c *taskController) List(w http.ResponseWriter, r *http.Request) error {

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

func (c *taskController) Create(w http.ResponseWriter, r *http.Request) error {

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

	taskId, warnings, err := c.service.CreateTask(
		r.Context(),
		projectId,
		userId,
		payload.Title,
		payload.Description,
		payload.Status,
		payload.Assignees)
	if err != nil {
		return fmt.Errorf("service create task: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[CreateTaskResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &CreateTaskResponse{
			Id:       taskId,
			Warnings: warnings,
		},
	})

	return nil
}

func (c *taskController) Get(w http.ResponseWriter, r *http.Request) error {

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

	assignees, err := c.service.GetTaskAssignees(
		r.Context(),
		projectId,
		taskId,
		userId,
	)
	if err != nil {
		return fmt.Errorf("service get task assignees: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[domain.TaskDetail]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &domain.TaskDetail{
			Task:      task,
			Assignees: assignees,
		},
	})

	return nil
}

func (c *taskController) Update(w http.ResponseWriter, r *http.Request) error {

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

	var payload UpdateTaskRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task accepts 'title', 'description', 'assignees' and 'status' only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("decode task payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Task 'title', 'description', 'assignees' and 'status' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("task payload validation: %w", err),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	err = c.service.UpdateTask(r.Context(),
		projectId,
		taskId,
		userId,
		payload.Title,
		payload.Description,
		payload.Status,
		payload.AssigneesToAdd,
		payload.AssigneesToRemove,
	)

	if errors.Is(err, apierr.ErrForbidden) {
		return &HTTPError{
			Code:    http.StatusForbidden,
			Message: "User is not allowed to update the task",
			ErrId:   ERR_ACCESS_DENIED,
			Err:     fmt.Errorf("service task update forbidden: %w", err),
		}
	} else if errors.Is(err, apierr.ErrInvalidValue) {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "User provided payload is invalid",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("service task update invalid: %w", err),
		}
	} else if err != nil {
		return fmt.Errorf("service task update: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Task updated successfully",
	})

	return nil
}

func (c *taskController) AddComment(w http.ResponseWriter, r *http.Request) error {

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

	var payload AddCommentRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Comment accepts 'user_id' and 'content' only",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("decode comment payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Comment 'user_id' or 'content' is missing",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("comment payload validation: %w", err),
		}
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	commentId, err := c.service.AddComment(
		r.Context(),
		projectId,
		taskId,
		userId,
		payload.Comment,
	)
	if err != nil {
		return fmt.Errorf("service add comment: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &commentId,
	})

	return nil
}

func (c *taskController) ListComments(w http.ResponseWriter, r *http.Request) error {

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

	comments, err := c.service.ListComments(
		r.Context(),
		projectId,
		taskId,
		userId,
	)
	if err != nil {
		return fmt.Errorf("service list comments: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedCommentsResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedCommentsResponse{
			Comments: comments,
			Page:     1,
			Limit:    len(comments),
			HasNext:  false,
		},
	})

	return nil
}

func (c *taskController) ListAssignedTasks(w http.ResponseWriter, r *http.Request) error {

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	tasks, err := c.service.AssignedTasks(r.Context(),
		userId)
	if err != nil {
		return fmt.Errorf("service assigned tasks: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedRecentTasksResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedRecentTasksResponse{
			Tasks: tasks,
		},
	})

	return nil
}

func (c *taskController) ListUnassignedTasks(w http.ResponseWriter, r *http.Request) error {

	userId, err := utils.GetUserId(r)
	if err != nil {
		return fmt.Errorf("get userId: %w", err)
	}

	tasks, err := c.service.UnassignedTasks(r.Context(),
		userId)
	if err != nil {
		return fmt.Errorf("service assigned tasks: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedRecentTasksResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedRecentTasksResponse{
			Tasks: tasks,
		},
	})

	return nil
}
