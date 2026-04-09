package models

import (
	"time"

	"github.com/ptracker/internal/domain"
)

type ProjectSummaryRow struct {
	ID              string    `gorm:"column:id"`
	Name            string    `gorm:"column:name"`
	Description     *string   `gorm:"column:description"`
	Skills          *string   `gorm:"column:skills"`
	OwnerID         string    `gorm:"column:owner_id"`
	UnassignedTasks int64     `gorm:"column:unassigned_tasks"`
	OngoingTasks    int64     `gorm:"column:ongoing_tasks"`
	CompletedTasks  int64     `gorm:"column:completed_tasks"`
	AbandonedTasks  int64     `gorm:"column:abandoned_tasks"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (r ProjectSummaryRow) ToProjectSummaryDomain() domain.ProjectSummary {
	return domain.ProjectSummary{
		Project: domain.Project{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Skills:      r.Skills,
			OwnerID:     r.OwnerID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		},
		UnassignedTasks: r.UnassignedTasks,
		OngoingTasks:    r.OngoingTasks,
		CompletedTasks:  r.CompletedTasks,
		AbandonedTasks:  r.AbandonedTasks,
	}
}

func MapToProjectSummaryDomain(ps []ProjectSummaryRow) []domain.ProjectSummary {
	d := []domain.ProjectSummary{}
	for _, r := range ps {
		d = append(d, r.ToProjectSummaryDomain())
	}

	return d
}

type ProjectPreviewRow struct {
	ID          string    `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Description *string   `gorm:"column:description"`
	Skills      *string   `gorm:"column:skills"`
	OwnerID     string    `gorm:"column:owner_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (r ProjectPreviewRow) ToProjectPreviewDomain() domain.ProjectPreview {
	return domain.ProjectPreview{
		Project: domain.Project{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Skills:      r.Skills,
			OwnerID:     r.OwnerID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		},
	}
}

func MapToProjectPreviewDomain(pv []ProjectPreviewRow) []domain.ProjectPreview {
	d := []domain.ProjectPreview{}
	for _, r := range pv {
		d = append(d, r.ToProjectPreviewDomain())
	}

	return d
}

type ProjectPublicDetailRow struct {
	ID               string    `gorm:"column:id"`
	Name             string    `gorm:"column:name"`
	Description      *string   `gorm:"column:description"`
	Skills           *string   `gorm:"column:skills"`
	OwnerID          string    `gorm:"column:owner_id"`
	OwnerUsername    string    `gorm:"column:owner_username"`
	OwnerDisplayName *string   `gorm:"column:owner_display_name"` // nullable
	OwnerEmail       string    `gorm:"column:owner_email"`
	OwnerAvatarURL   *string   `gorm:"column:owner_avatar_url"`
	MemberCount      int64     `gorm:"column:member_count"`
	UnassignedTasks  int64     `gorm:"column:unassigned_tasks"`
	OngoingTasks     int64     `gorm:"column:ongoing_tasks"`
	CompletedTasks   int64     `gorm:"column:completed_tasks"`
	AbandonedTasks   int64     `gorm:"column:abandoned_tasks"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (r ProjectPublicDetailRow) ToProjectPublicDetailDomain() domain.ProjectPublicDetail {
	return domain.ProjectPublicDetail{
		ProjectSummary: domain.ProjectSummary{
			Project: domain.Project{
				ID:          r.ID,
				Name:        r.Name,
				Description: r.Description,
				Skills:      r.Skills,
				OwnerID:     r.OwnerID,
				CreatedAt:   r.CreatedAt,
				UpdatedAt:   r.UpdatedAt,
			},
			UnassignedTasks: r.UnassignedTasks,
			OngoingTasks:    r.OngoingTasks,
			CompletedTasks:  r.CompletedTasks,
			AbandonedTasks:  r.AbandonedTasks,
		},
		MemberCount: r.MemberCount,
		Owner: domain.Avatar{
			UserID:      r.OwnerID,
			Username:    r.OwnerUsername,
			DisplayName: r.OwnerDisplayName,
			Email:       r.OwnerEmail,
			AvatarURL:   r.OwnerAvatarURL,
		},
	}
}

type AssigneeRow struct {
	ProjectID string    `gorm:"column:project_id"`
	TaskID    string    `gorm:"column:task_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	AssigneeID          string  `gorm:"column:assignee_id"`
	AssigneeUsername    string  `gorm:"column:assignee_username"`
	AssigneeDisplayName *string `gorm:"column:assignee_display_name"`
	AssigneeEmail       string  `gorm:"column:assignee_email"`
	AssigneeAvatarURL   *string `gorm:"column:assignee_avatar_url"`
}

func (a AssigneeRow) ToAssigneeDomain() domain.Assignee {
	return domain.Assignee{
		ProjectID: a.ProjectID,
		TaskID:    a.TaskID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Avatar: domain.Avatar{
			UserID:      a.AssigneeID,
			Username:    a.AssigneeUsername,
			DisplayName: a.AssigneeDisplayName,
			Email:       a.AssigneeEmail,
			AvatarURL:   a.AssigneeAvatarURL,
		},
	}
}

type ProjectTaskItemRow struct {
	ID          string    `gorm:"column:id"`
	Title       string    `gorm:"column:title"`
	Description *string   `gorm:"column:description"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	ProjectID string        `gorm:"column:project_id"`
	Assignees []AssigneeRow `gorm:"column:assignees"`
}

func (t ProjectTaskItemRow) ToProjectTaskItemDomain() domain.ProjectTaskItem {
	item := domain.ProjectTaskItem{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Assignees:   []domain.Assignee{},
	}

	for _, assignee := range t.Assignees {
		item.Assignees = append(item.Assignees, assignee.ToAssigneeDomain())
	}

	return item
}

func MapToProjectTaskItemDomain(ts []ProjectTaskItemRow) []domain.ProjectTaskItem {
	d := []domain.ProjectTaskItem{}
	for _, r := range ts {
		d = append(d, r.ToProjectTaskItemDomain())
	}

	return d
}

type MembershipRow struct {
	ProjectID string    `gorm:"column:project_id"`
	Role      string    `gorm:"column:role"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	UserID      string  `gorm:"column:user_id"`
	Username    string  `gorm:"column:username"`
	DisplayName *string `gorm:"column:display_name"`
	Email       string  `gorm:"column:email"`
	AvatarURL   *string `gorm:"column:avatar_url"`
}

func (m MembershipRow) ToMembershipDomain() domain.Membership {
	return domain.Membership{
		ProjectID: m.ProjectID,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Avatar: domain.Avatar{
			UserID:      m.UserID,
			Username:    m.Username,
			Email:       m.Email,
			DisplayName: m.DisplayName,
			AvatarURL:   m.AvatarURL,
		},
	}
}

func MapToMembershipDomain(ms []MembershipRow) []domain.Membership {
	d := []domain.Membership{}
	for _, r := range ms {
		d = append(d, r.ToMembershipDomain())
	}
	return d
}

type JoinRequestRow struct {
	ProjectID string    `gorm:"column:project_id"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	UserID      string  `gorm:"column:user_id"`
	Username    string  `gorm:"column:username"`
	DisplayName *string `gorm:"column:display_name"`
	Email       string  `gorm:"column:email"`
	AvatarURL   *string `gorm:"column:avatar_url"`
}

func (j JoinRequestRow) ToJoinRequestDomain() domain.JoinRequest {
	return domain.JoinRequest{
		ProjectID: j.ProjectID,
		Status:    j.Status,
		CreatedAt: j.CreatedAt,
		UpdatedAt: j.UpdatedAt,
		Avatar: domain.Avatar{
			UserID:      j.UserID,
			Username:    j.Username,
			DisplayName: j.DisplayName,
			Email:       j.Email,
			AvatarURL:   j.AvatarURL,
		},
	}
}

func MapToJoinRequestDomain(js []JoinRequestRow) []domain.JoinRequest {
	d := []domain.JoinRequest{}
	for _, r := range js {
		d = append(d, r.ToJoinRequestDomain())
	}
	return d
}

type CommentRow struct {
	ID        string    `gorm:"column:id"`
	ProjectID string    `gorm:"column:project_id"`
	TaskID    string    `gorm:"column:task_id"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	UserID      string  `gorm:"column:user_id"`
	Username    string  `gorm:"column:username"`
	DisplayName *string `gorm:"column:display_name"`
	Email       string  `gorm:"column:email"`
	AvatarURL   *string `gorm:"column:avatar_url"`
}

func (c CommentRow) ToCommentDomain() domain.Comment {
	return domain.Comment{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		TaskID:    c.TaskID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Avatar: domain.Avatar{
			UserID:      c.UserID,
			Username:    c.Username,
			DisplayName: c.DisplayName,
			Email:       c.Email,
			AvatarURL:   c.AvatarURL,
		},
	}
}

func MapToCommentDomain(cs []CommentRow) []domain.Comment {
	d := []domain.Comment{}
	for _, r := range cs {
		d = append(d, r.ToCommentDomain())
	}
	return d
}

type DashboardTaskItemRow struct {
	ID          string    `gorm:"column:id"`
	Title       string    `gorm:"column:title"`
	Description *string   `gorm:"column:description"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`

	ProjectID   string `gorm:"column:project_id"`
	ProjectName string `gorm:"column:project_name"`
}

func (d DashboardTaskItemRow) ToDashboardTaskItemDomain() domain.DashboardTaskItem {
	return domain.DashboardTaskItem{
		ID:          d.ID,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		ProjectID:   d.ProjectID,
		ProjectName: d.ProjectName,
	}
}

func MapToDashboardTaskItemDomain(ts []DashboardTaskItemRow) []domain.DashboardTaskItem {
	d := []domain.DashboardTaskItem{}
	for _, r := range ts {
		d = append(d, r.ToDashboardTaskItemDomain())
	}

	return d
}
