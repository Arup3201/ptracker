package services

import (
	"context"
	"sync"

	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories"
	"gorm.io/gorm"
)

type storage struct {
	mu sync.Mutex

	db *gorm.DB

	sessionRepo     interfaces.SessionRepository
	userRepo        interfaces.UserRepository
	projectRepo     interfaces.ProjectRepository
	taskRepo        interfaces.TaskRepository
	commentRepo     interfaces.CommentRepository
	membershipRepo  interfaces.MembershipRepository
	assigneeRepo    interfaces.AssigneeRepository
	listRepo        interfaces.ListRepository
	joinRequestRepo interfaces.JoinRequestRepository
	publicRepo      interfaces.PublicRepository

	inMemory    interfaces.InMemory
	rateLimiter interfaces.RateLimiter
}

func NewStorage(db *gorm.DB,
	memory interfaces.InMemory,
	rateLimiter interfaces.RateLimiter) interfaces.Store {
	s := &storage{}
	s.db = db
	s.sessionRepo = repositories.NewSessionRepo(db)
	s.userRepo = repositories.NewUserRepo(db)
	s.projectRepo = repositories.NewProjectRepo(db)
	s.taskRepo = repositories.NewTaskRepo(db)
	s.commentRepo = repositories.NewCommentRepo(db)
	s.membershipRepo = repositories.NewMembershipRepo(db)
	s.assigneeRepo = repositories.NewAssigneeRepo(db)
	s.listRepo = repositories.NewListRepo(db)
	s.joinRequestRepo = repositories.NewJoinRequestRepo(db)
	s.publicRepo = repositories.NewPublicRepo(db)

	s.inMemory = memory
	s.rateLimiter = rateLimiter

	return s
}

func (s *storage) WithTx(ctx context.Context, fn func(txStore interfaces.Store) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewStorage(tx, s.inMemory, s.rateLimiter))
	})

	return err
}

func (s *storage) Session() interfaces.SessionRepository {
	return s.sessionRepo
}

func (s *storage) User() interfaces.UserRepository {
	return s.userRepo
}

func (s *storage) Project() interfaces.ProjectRepository {
	return s.projectRepo
}

func (s *storage) Task() interfaces.TaskRepository {
	return s.taskRepo
}

func (s *storage) Comment() interfaces.CommentRepository {
	return s.commentRepo
}

func (s *storage) Membership() interfaces.MembershipRepository {
	return s.membershipRepo
}

func (s *storage) Assignee() interfaces.AssigneeRepository {
	return s.assigneeRepo
}

func (s *storage) List() interfaces.ListRepository {
	return s.listRepo
}

func (s *storage) JoinRequest() interfaces.JoinRequestRepository {
	return s.joinRequestRepo
}

func (s *storage) Public() interfaces.PublicRepository {
	return s.publicRepo
}

func (s *storage) InMemory() interfaces.InMemory {
	return s.inMemory
}

func (s *storage) RateLimiter() interfaces.RateLimiter {
	return s.rateLimiter
}
