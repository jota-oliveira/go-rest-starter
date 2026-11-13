package items

import (
	"context"
	"errors"
	"time"

	"example.com/go-rest-starter/internal/domain/item"
	"github.com/google/uuid"
)

var ErrNotAuthenticated = errors.New("not authenticated")

type Repository interface {
	Save(context.Context, item.Item) error
	ListByOwner(context.Context, string) ([]item.Item, error)
}

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, owner, name string) (item.Item, error) {
	if owner == "" {
		return item.Item{}, ErrNotAuthenticated
	}
	clean, err := item.ValidateName(name)
	if err != nil {
		return item.Item{}, err
	}
	obj := item.Item{ID: uuid.NewString(), OwnerID: owner, Name: clean, CreatedAt: time.Now().UTC()}
	if err := s.repo.Save(ctx, obj); err != nil {
		return item.Item{}, err
	}
	return obj, nil
}
func (s *Service) List(ctx context.Context, owner string) ([]item.Item, error) {
	if owner == "" {
		return nil, ErrNotAuthenticated
	}
	return s.repo.ListByOwner(ctx, owner)
}
