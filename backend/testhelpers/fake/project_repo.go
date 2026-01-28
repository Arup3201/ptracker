package fake

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
)

type ProjectRepo struct {
	mu       sync.RWMutex
	projects map[string]*domain.Project
	nextID   int
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{
		projects: make(map[string]*domain.Project),
		nextID:   1,
	}
}

func (r *ProjectRepo) Create(
	ctx context.Context,
	title string,
	description, skills *string,
	owner string,
) (string, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.nextID)
	r.nextID++

	now := time.Now()
	p := &domain.Project{
		Id:          id,
		Name:        title,
		Description: description,
		Skills:      skills,
		Owner:       owner,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	r.projects[id] = p
	return id, nil
}

func (r *ProjectRepo) Get(ctx context.Context, id string) (*domain.ListedProject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.projects[id]
	if !ok {
		return nil, apierr.ErrResourceNotFound
	}

	cp := &domain.ListedProject{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Skills:      p.Skills,
	}
	return cp, nil
}

func (r *ProjectRepo) All(ctx context.Context, userId string) ([]*domain.ListedProject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.ListedProject
	for _, p := range r.projects {
		if p.Owner == userId {
			result = append(result, &domain.ListedProject{
				Id:              p.Id,
				Name:            p.Name,
				Description:     p.Description,
				Skills:          p.Skills,
				UnassignedTasks: 0,
				OngoingTasks:    0,
				CompletedTasks:  0,
				AbandonedTasks:  0,
				CreatedAt:       p.CreatedAt,
				UpdatedAt:       p.UpdatedAt,
			})
		}
	}
	return result, nil
}

func (r *ProjectRepo) WithProject(p domain.Project) *ProjectRepo {
	r.mu.Lock()
	defer r.mu.Unlock()

	if p.Id == "" {
		p.Id = strconv.Itoa(r.nextID)
		r.nextID++
	}

	r.projects[p.Id] = &p
	return r
}
