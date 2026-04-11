package assignees

import (
	"context"
	"fmt"
	"time"

	"github.com/ptracker/core"
	"github.com/ptracker/core/members"
)

type Assignee struct {
	ProjectID   string    `json:"project_id"`
	TaskID      string    `json:"task_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	core.Avatar `json:"avatar"`
}

type AssigneeService struct {
	memberRepo   *members.MemberRepository
	assigneeRepo *AssigneeRepository
}

func NewAssigneeService() *AssigneeService {
	return &AssigneeService{}
}

func (s *AssigneeService) AddAssignee(ctx context.Context,
	projectID, taskID, userID, assigneeID string) error {

	var err error

	err = core.NeedsToBeAnOwner(ctx, s.memberRepo, projectID, userID)
	if err != nil {
		return fmt.Errorf("needs to be a owner: %w", err)
	}

	if err = s.assigneeRepo.Is(ctx, projectID, taskID, userID); err != core.ErrNotFound {
		return core.ErrDuplicate
	}

	err = s.assigneeRepo.Create(ctx, projectID, taskID, assigneeID)
	if err != nil {
		return fmt.Errorf("assignee repository create: %w", err)
	}

	return nil
}

func (s *AssigneeService) RemoveAssignee(ctx context.Context,
	projectID, taskID, userID, assigneeID string) error {

	var err error

	err = core.NeedsToBeAnOwner(ctx, s.memberRepo, projectID, userID)
	if err != nil {
		return fmt.Errorf("needs to be a owner: %w", err)
	}

	if err = s.assigneeRepo.Is(ctx, projectID, taskID, userID); err == core.ErrNotFound {
		return core.ErrNotFound
	}

	err = s.assigneeRepo.Delete(ctx, projectID, taskID, assigneeID)
	if err != nil {
		return fmt.Errorf("assignee repository create: %w", err)
	}

	return nil
}
