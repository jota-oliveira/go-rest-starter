package items_test

import (
	"context"
	"example.com/go-rest-starter/internal/adapters/persistence/memory"
	"example.com/go-rest-starter/internal/application/items"
	"testing"
)

func TestCreateAndListByOwner(t *testing.T) {
	s := items.NewService(memory.New())
	ctx := context.Background()
	_, err := s.Create(ctx, "alice", "  example  ")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(ctx, "bob", "other"); err != nil {
		t.Fatal(err)
	}
	list, err := s.List(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "example" {
		t.Fatalf("unexpected list: %+v", list)
	}
	if _, err = s.Create(ctx, "alice", ""); err == nil {
		t.Fatal("expected validation error")
	}
}
