package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ptracker/core"
	"github.com/ptracker/core/members"
	"github.com/ptracker/core/projects"
	"github.com/ptracker/core/tasks"
	"github.com/ptracker/core/users"
)

const (
	NT_TASK_ADDED       = "task_added"
	NT_TASK_UPDATED     = "task_updated"
	NT_ASSIGNEE_ADDED   = "assignee_added"
	NT_ASSIGNEE_REMOVED = "assignee_removed"
	NT_JOIN_REQUESTED   = "join_requested"
	NT_JOIN_RESPONDED   = "join_responded"
	NT_COMMENT_ADDED    = "comment_added"
)

type ProjectBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TaskBody struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskUpdateBody struct {
	To    string `json:"to"`
	Field string `json:"field"`
}

type TaskAdded struct {
	Project ProjectBody `json:"project"`
	Task    TaskBody    `json:"task"`
}

type TaskUpdated struct {
	Project ProjectBody      `json:"project"`
	Task    TaskBody         `json:"task"`
	Updates []TaskUpdateBody `json:"updates"`
	Updater core.Avatar      `json:"updater"`
}

type AssigneeUpdated struct {
	Project  ProjectBody `json:"project"`
	Task     TaskBody    `json:"task"`
	Assignee core.Avatar `json:"assignee"`
}

type JoinRequested struct {
	Project   ProjectBody `json:"project"`
	Requestor core.Avatar `json:"requestor"`
}

type JoinResponded struct {
	Project   ProjectBody `json:"project"`
	Responder core.Avatar `json:"responder"`
	Status    string      `json:"status"`
}

type CommentAdded struct {
	Project   ProjectBody `json:"project"`
	Task      TaskBody    `json:"task"`
	Commenter core.Avatar `json:"commenter"`
}

type Notification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Body   any    `json:"body"`
}

type NotificationService struct {
	projectRepo      *projects.ProjectRepository
	taskRepo         *tasks.TaskRepository
	membershipRepo   *members.MemberRepository
	userRepo         *users.UserRepository
	notificationRepo *NotificationRepository
}

func NewNotificationService(
	projectRepo *projects.ProjectRepository,
	taskRepo *tasks.TaskRepository,
	membershipRepo *members.MemberRepository,
	userRepo *users.UserRepository,
	notificationRepo *NotificationRepository,
) *NotificationService {
	return &NotificationService{
		projectRepo:      projectRepo,
		taskRepo:         taskRepo,
		membershipRepo:   membershipRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) TaskAdded(ctx context.Context,
	projectID, taskID string) error {

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project repository Get: %w", err)
	}

	task, err := s.taskRepo.Get(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task repository Get: %w", err)
	}

	members, err := s.membershipRepo.List(ctx, projectID)
	if err != nil {
		return fmt.Errorf("membership repository Get: %w", err)
	}

	body, _ := json.Marshal(TaskAdded{
		Project: ProjectBody{
			ID:   projectID,
			Name: project.Name,
		},
		Task: TaskBody{
			ID:    taskID,
			Title: task.Title,
		},
	})
	for _, m := range members {
		if m.Role.String == core.ROLE_OWNER {
			continue
		}

		_, err := s.notificationRepo.Create(ctx, m.UserID, NT_TASK_ADDED, body, false)
		if err != nil {
			return fmt.Errorf("notification repository Create: %w", err)
		}
	}

	return nil
}

func (s *NotificationService) TaskUpdated(ctx context.Context,
	projectID, taskID string,
	title, description, status *string,
	updater string) error {

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project repository Get: %w", err)
	}

	task, err := s.taskRepo.Get(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task repository Get: %w", err)
	}

	updaterInfo, err := s.userRepo.Get(ctx, updater)
	if err != nil {
		return fmt.Errorf("user repository Get: %w", err)
	}

	updates := []TaskUpdateBody{}
	if title != nil {
		updates = append(updates, TaskUpdateBody{
			To:    *title,
			Field: "Title",
		})
	}
	if description != nil {
		updates = append(updates, TaskUpdateBody{
			To:    *description,
			Field: "Description",
		})
	}
	if status != nil {
		updates = append(updates, TaskUpdateBody{
			To:    *status,
			Field: "Status",
		})
	}

	body, _ := json.Marshal(TaskUpdated{
		Project: ProjectBody{
			ID:   projectID,
			Name: project.Name,
		},
		Task: TaskBody{
			ID:    taskID,
			Title: task.Title,
		},
		Updates: updates,
		Updater: core.Avatar{
			UserID:      updaterInfo.ID,
			Username:    updaterInfo.Username,
			DisplayName: updaterInfo.DisplayName,
			Email:       updaterInfo.Email,
			AvatarURL:   updaterInfo.AvatarURL,
		},
	})

	if project.OwnerID != updater {
		_, err := s.notificationRepo.Create(
			ctx,
			project.OwnerID,
			NT_TASK_UPDATED,
			body,
			false,
		)
		if err != nil {
			return fmt.Errorf("notification repository Create: %w", err)
		}
	}

	for _, assignee := range task.Assignees {
		if assignee.AssigneeID == updater {
			continue
		}

		_, err := s.notificationRepo.Create(
			ctx,
			assignee.AssigneeID,
			NT_TASK_UPDATED,
			body,
			false,
		)
		if err != nil {
			return fmt.Errorf("notification repository Create: %w", err)
		}
	}

	return nil
}
