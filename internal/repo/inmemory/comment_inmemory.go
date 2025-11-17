package inmemory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
)

var ErrNotFound = errors.New("comment not found")

type InMemoryCommentRepo struct {
	mu       sync.RWMutex
	nextID   int64
	byID     map[int64]*entity.Comment
	byEntity map[int64][]*entity.Comment // entityID -> slice of comments
}

func NewInMemoryCommentRepo() *InMemoryCommentRepo {
	return &InMemoryCommentRepo{
		nextID:   1,
		byID:     make(map[int64]*entity.Comment),
		byEntity: make(map[int64][]*entity.Comment),
	}
}

func (r *InMemoryCommentRepo) Create(ctx context.Context, c *entity.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = r.nextID
	r.nextID++
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
		c.UpdatedAt = c.CreatedAt
	}
	// store copy to avoid external mutation
	cc := *c
	r.byID[c.ID] = &cc
	r.byEntity[c.EntityID] = append(r.byEntity[c.EntityID], &cc)
	return nil
}

func (r *InMemoryCommentRepo) Update(ctx context.Context, c *entity.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.byID[c.ID]
	if !ok {
		return ErrNotFound
	}
	existing.Text = c.Text
	existing.UpdatedAt = time.Now()
	// also update in slice
	list := r.byEntity[c.EntityID]
	for i := range list {
		if list[i].ID == c.ID {
			list[i] = existing
			break
		}
	}
	return nil
}

func (r *InMemoryCommentRepo) GetByID(ctx context.Context, id int64) (*entity.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	// return copy
	copy := *c
	return &copy, nil
}

func (r *InMemoryCommentRepo) ListByEntity(ctx context.Context, entityID int64, page, limit int, sortAsc bool) ([]*entity.Comment, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, ok := r.byEntity[entityID]
	if !ok {
		return []*entity.Comment{}, 0, nil
	}

	// create copy slice
	copyList := make([]*entity.Comment, len(list))
	for i := range list {
		v := *list[i]
		copyList[i] = &v
	}

	// sort by CreatedAt
	sort.Slice(copyList, func(i, j int) bool {
		if sortAsc {
			return copyList[i].CreatedAt.Before(copyList[j].CreatedAt)
		}
		return copyList[i].CreatedAt.After(copyList[j].CreatedAt)
	})

	total := int64(len(copyList))
	if limit <= 0 {
		return copyList, total, nil
	}
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}
	if start >= len(copyList) {
		return []*entity.Comment{}, total, nil
	}
	end := start + limit
	if end > len(copyList) {
		end = len(copyList)
	}
	return copyList[start:end], total, nil
}
