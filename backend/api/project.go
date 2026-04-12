package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ptracker/core"
	"github.com/ptracker/core/members"
	"github.com/ptracker/core/projects"
	"github.com/ptracker/core/requests"
	"github.com/ptracker/core/users"
)

type CreateProjectRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Skills      *string `json:"skills"`
}

type ProjectDetail struct {
	projects.ProjectSummary
	Role        string      `json:"role"`
	MemberCount int64       `json:"members_count"`
	Owner       core.Avatar `json:"owner"`
}

type ListedProjectSummaries struct {
	Projects []projects.ProjectSummary `json:"projects"`
}

type ListedJoinRequests struct {
	JoinRequests []requests.JoinRequest `json:"join_requests"`
}

type ListedMembers struct {
	Members []members.Member `json:"members"`
}

type UpdateJoinRequest struct {
	UserID     string `json:"user_id" validate:"required"`
	JoinStatus string `json:"join_status" validate:"required"`
}

type ProjectApi struct {
	projectService     *projects.ProjectService
	userService        *users.UserService
	memberService      *members.MemberService
	joinRequestService *requests.JoinRequestService
}

func (api *ProjectApi) Create(w http.ResponseWriter, r *http.Request) error {

	var payload CreateProjectRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project can only have 'name', 'description' and 'skills' fields.",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("decode payload: %w", err),
		}
	}
	if err := validator.New().Struct(payload); err != nil {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'name' is required.",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("validate payload: %w", err),
		}
	}

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get user Id: %w", err)
	}

	projectID, err := api.projectService.Create(r.Context(),
		payload.Name, payload.Description, payload.Skills, userID)
	if err != nil {
		return fmt.Errorf("store create project: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &projectID,
	})
	return nil
}

func (api *ProjectApi) Get(w http.ResponseWriter, r *http.Request) error {

	projectID := r.PathValue("id")
	if projectID == "" {
		return core.ErrInvalidValue
	}

	projectSummary, err := api.projectService.Get(r.Context(), projectID)
	if err != nil {
		return fmt.Errorf("database get project details: %w", err)
	}

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get context user: %w", err)
	}

	role, err := api.memberService.GetRole(r.Context(), projectID, userID)
	if !errors.Is(err, core.ErrNotFound) && err != nil {
		return fmt.Errorf("member service get role: %w", err)
	}

	if role == core.ROLE_MEMBER || role == core.ROLE_OWNER {
		memberCount, err := api.memberService.Count(r.Context(), projectID, userID)
		if err != nil {
			return fmt.Errorf("member service count: %w", err)
		}

		owner, err := api.userService.Get(r.Context(), userID)
		if err != nil {
			return fmt.Errorf("user service get: %w", err)
		}

		json.NewEncoder(w).Encode(HTTPSuccessResponse[ProjectDetail]{
			Status: RESPONSE_SUCCESS_STATUS,
			Data: &ProjectDetail{
				ProjectSummary: *projectSummary,
				MemberCount:    memberCount,
				Owner: core.Avatar{
					UserID:      owner.ID,
					Username:    owner.Username,
					Email:       owner.Email,
					DisplayName: owner.DisplayName,
					AvatarURL:   owner.AvatarURL,
				},
			},
		})
	} else {
		json.NewEncoder(w).Encode(HTTPSuccessResponse[projects.ProjectSummary]{
			Status: RESPONSE_SUCCESS_STATUS,
			Data:   projectSummary,
		})
	}

	return nil
}

func (api *ProjectApi) ListRecentlyCreated(w http.ResponseWriter, r *http.Request) error {

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get userID: %w", err)
	}

	projects, err := api.projectService.RecentlyCreated(r.Context(), userID)
	if err != nil {
		return fmt.Errorf("project service recently created: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedProjectSummaries]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedProjectSummaries{
			Projects: projects,
		},
	})

	return nil
}

func (api *ProjectApi) ListRecentlyJoined(w http.ResponseWriter, r *http.Request) error {
	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get userID: %w", err)
	}

	projects, err := api.projectService.RecentlyJoined(r.Context(), userID)
	if err != nil {
		return fmt.Errorf("project service recently joined: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedProjectSummaries]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedProjectSummaries{
			Projects: projects,
		},
	})

	return nil
}

func (api *ProjectApi) ListMembers(w http.ResponseWriter, r *http.Request) error {

	projectID := r.PathValue("id")
	if projectID == "" {
		return core.ErrInvalidValue
	}

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get userID: %w", err)
	}

	members, err := api.memberService.AllMembers(r.Context(),
		projectID,
		userID)
	if err != nil {
		return fmt.Errorf("member service all members: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedMembers]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedMembers{
			Members: members,
		},
	})

	return nil
}

func (api *ProjectApi) ListJoinRequests(w http.ResponseWriter, r *http.Request) error {

	projectID := r.PathValue("id")
	if projectID == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Project 'id' can't be empty",
			ErrId:   ERR_INVALID_BODY,
			Err:     fmt.Errorf("empty 'id' provided"),
		}
	}

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get userID: %w", err)
	}

	joinRequests, err := api.joinRequestService.List(r.Context(),
		projectID,
		userID)
	if err != nil {
		return fmt.Errorf("join request service list: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedJoinRequests]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &ListedJoinRequests{
			JoinRequests: joinRequests,
		},
	})

	return nil
}

func (api *ProjectApi) RespondToJoinRequest(w http.ResponseWriter, r *http.Request) error {

	projectId := r.PathValue("id")
	if projectId == "" {
		return core.ErrInvalidValue
	}

	var payload UpdateJoinRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		return core.ErrInvalidValue
	}
	if err := validator.New().Struct(payload); err != nil {
		return core.ErrInvalidValue
	}

	if payload.UserID == "" {
		return core.ErrInvalidValue
	}
	if payload.JoinStatus == "" {
		return core.ErrInvalidValue
	}

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get projects userID: %w", err)
	}

	err = api.joinRequestService.Respond(
		r.Context(),
		projectId,
		userID,
		payload.UserID,
		payload.JoinStatus,
	)
	if err != nil {
		return fmt.Errorf("join request service respond: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Join request status updated",
	})

	return nil
}
