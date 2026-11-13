package memory

import (
	"context"
	"example.com/go-rest-starter/internal/domain/item"
	"sort"
	"sync"
)

type Repository struct {
	mu    sync.RWMutex
	items map[string]item.Item
}

func New() *Repository { return &Repository{items: make(map[string]item.Item)} }
func (r *Repository) Save(_ context.Context, obj item.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[obj.ID] = obj
	return nil
}
func (r *Repository) ListByOwner(_ context.Context, owner string) ([]item.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]item.Item, 0)
	for _, v := range r.items {
		if v.OwnerID == owner {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}
