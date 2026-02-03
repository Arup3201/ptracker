package interfaces

import (
	"context"
)

type Store interface {
	WithTx(ctx context.Context, fn func(txStore Store) error) error

	Session() SessionRepository
	User() UserRepository
	Project() ProjectRepository
	Task() TaskRepository
	Comment() CommentRepository
	Role() RoleRepository
	Assignee() AssigneeRepository
	List() ListRepository
	JoinRequest() JoinRequestRepository
	Public() PublicRepository

	InMemory() InMemory
	RateLimiter() RateLimiter
}
