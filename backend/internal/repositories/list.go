package repositories

import (
	"context"
	"fmt"

	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type ListRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) interfaces.ListRepository {
	return &ListRepo{
		db: db,
	}
}

func (r *ListRepo) PrivateProjects(ctx context.Context, userId string) ([]domain.ProjectSummary, error) {

	var rows []models.ProjectSummaryRow
	err := r.db.WithContext(ctx).
		Table("projects p").
		Select(`p.id, p.name, p.description, p.skills, p.owner_id, 
				ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, 
				p.created_at, p.updated_at`).
		Joins("INNER JOIN memberships as m ON m.project_id=p.id").
		Joins("LEFT JOIN project_summary as ps ON ps.id=p.id").
		Where("m.user_id = ?", userId).
		Scan(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToProjectSummaryDomain(rows), nil
}

func (r *ListRepo) PublicProjects(ctx context.Context, userId string) ([]domain.ProjectPreview, error) {

	var rows []models.ProjectPreviewRow
	err := r.db.WithContext(ctx).
		Table("projects").
		Select(`id, name, description, skills, owner_id, 
			created_at, updated_at`).
		Where("owner_id != ? AND NOT EXISTS (?)",
			userId,
			r.db.WithContext(ctx).
				Table("memberships").
				Select("1").
				Where("project_id = projects.id AND user_id = ?", userId),
		).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToProjectPreviewDomain(rows), nil
}

func (r *ListRepo) Tasks(ctx context.Context,
	projectId string) ([]domain.ProjectTaskItem, error) {

	query := `SELECT 
		t.id, t.title, t.status, t.created_at, t.updated_at, 
		COALESCE(
			json_agg(
			json_build_object(
				'project_id', a.project_id,
				'task_id', a.task_id,
				'assignee_id', a.user_id,
				'assignee_username', u.username,
				'assignee_display_name', u.display_name,
				'assignee_email', u.email,
				'assignee_avatar_url', u.avatar_url
			)
			) FILTER (WHERE a.user_id IS NOT NULL),
			'[]'
		) AS assignees 
		FROM tasks AS t 
		LEFT JOIN assignees AS a ON a.task_id=t.id 
		LEFT JOIN users AS u ON u.id=a.user_id 
		WHERE t.project_id=? 
		GROUP BY t.id`

	var rows []models.ProjectTaskItemRow
	err := r.db.Raw(query, projectId).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("gorm db raw scan: %w", err)
	}

	return models.MapToProjectTaskItemDomain(rows), nil
}

func (r *ListRepo) Members(ctx context.Context,
	projectId string) ([]domain.Membership, error) {

	var rows []models.MembershipRow
	err := r.db.WithContext(ctx).
		Table("memberships m").
		Select(`m.project_id, m.role, m.created_at, m.updated_at, 
				u.id as member_user_id, 
				u.username as member_username, 
				u.display_name as member_display_name, 
				u.email as member_email, 
				u.avatar_url as member_avatar_url`).
		Joins("INNER JOIN users AS u ON m.user_id=u.id").
		Where("m.project_id = ?", projectId).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToMembershipDomain(rows), nil
}

func (r *ListRepo) JoinRequests(ctx context.Context, projectId string) ([]domain.JoinRequest, error) {

	var rows []models.JoinRequestRow
	err := r.db.WithContext(ctx).
		Table("join_requests jr").
		Select(`jr.project_id, jr.status, jr.created_at, jr.updated_at, 
				u.id as requestor_user_id, 
				u.username as requestor_username, 
				u.display_name as requestor_diplay_name, 
				u.email as requestor_email, 
				u.avatar_url as requestor_avatar_url`).
		Joins("INNER JOIN users as u ON u.id=jr.user_id").
		Where("jr.project_id = ?", projectId).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}

	return models.MapToJoinRequestDomain(rows), nil
}

func (r *ListRepo) Comments(ctx context.Context,
	projectId, taskId string) ([]domain.Comment, error) {

	var rows []models.CommentRow
	err := r.db.WithContext(ctx).
		Table("comments c").
		Select(`c.id, c.project_id, c.task_id, c.content, 
				c.created_at, c.updated_at, 
				u.id, u.username, u.display_name, u.email, u.avatar_url`).
		Joins("INNER JOIN users as u ON u.id=c.user_id").
		Where("c.project_id = ? AND c.task_id = ?", projectId, taskId).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("db query context: %w", err)
	}

	return models.MapToCommentDomain(rows), nil
}

func (r *ListRepo) RecentlyCreatedProjects(ctx context.Context,
	userId string, n int) ([]domain.ProjectSummary, error) {

	var rows []models.ProjectSummaryRow
	err := r.db.WithContext(ctx).
		Table("projects p").
		Select(`p.id, p.name, p.description, p.skills, p.owner_id, 
				ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks, 
				p.created_at, p.updated_at`).
		Joins("LEFT JOIN project_summary as ps ON ps.id=p.id").
		Where("p.owner_id = ?", userId).
		Order("p.created_at DESC").
		Limit(n).
		Scan(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToProjectSummaryDomain(rows), nil
}

func (r *ListRepo) RecentlyJoinedProjects(ctx context.Context,
	userId string,
	n int) ([]domain.ProjectSummary, error) {

	var rows []models.ProjectSummaryRow
	err := r.db.WithContext(ctx).
		Table("projects p").
		Select(`p.id, p.name, p.description, p.skills, p.owner_id,
				ps.unassigned_tasks, ps.ongoing_tasks, ps.completed_tasks, ps.abandoned_tasks,
				p.created_at, p.updated_at`).
		Joins("INNER JOIN memberships as m ON m.project_id=p.id").
		Joins("LEFT JOIN project_summary as ps ON ps.id=p.id").
		Where("m.user_id = ? AND m.role = ?", userId, domain.ROLE_MEMBER).
		Order("m.created_at DESC").
		Limit(n).
		Scan(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToProjectSummaryDomain(rows), nil
}

func (r *ListRepo) RecentlyAssignedTasks(ctx context.Context,
	userId string,
	n int) ([]domain.DashboardTaskItem, error) {

	var rows []models.DashboardTaskItemRow
	err := r.db.WithContext(ctx).
		Table("tasks t").
		Select(`t.id, t.project_id, t.title, t.status, t.created_at, t.updated_at,
				p.name as project_name`).
		Joins("INNER JOIN projects AS p ON t.project_id=p.id").
		Joins("INNER JOIN assignees AS a ON a.task_id=t.id").
		Where("a.user_id = ?", userId).
		Order("a.created_at DESC").
		Limit(n).
		Scan(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToDashboardTaskItemDomain(rows), nil
}

func (r *ListRepo) RecentlyUnassignedTasks(ctx context.Context,
	userId string,
	n int) ([]domain.DashboardTaskItem, error) {

	var rows []models.DashboardTaskItemRow
	err := r.db.WithContext(ctx).
		Table("tasks t").
		Select(`t.id, t.project_id, t.title, t.status, t.created_at, t.updated_at,
					p.name as project_name`).
		Joins("INNER JOIN projects AS p ON t.project_id=p.id").
		Where("p.owner_id = ? AND t.status = ?", userId, domain.TASK_STATUS_UNASSIGNED).
		Order("t.created_at DESC").
		Limit(n).
		Scan(&rows).
		Error
	if err != nil {
		return nil, fmt.Errorf("gorm db scan: %w", err)
	}

	return models.MapToDashboardTaskItemDomain(rows), nil
}
