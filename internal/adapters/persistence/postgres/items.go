package postgres

import (
	"context"
	"example.com/go-rest-starter/internal/domain/item"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }
func (r *Repository) Save(ctx context.Context, obj item.Item) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO items (id,owner_id,name,created_at) VALUES ($1,$2,$3,$4)", obj.ID, obj.OwnerID, obj.Name, obj.CreatedAt)
	return err
}
func (r *Repository) ListByOwner(ctx context.Context, owner string) ([]item.Item, error) {
	rows, err := r.pool.Query(ctx, "SELECT id::text,owner_id,name,created_at FROM items WHERE owner_id=$1 ORDER BY created_at DESC", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]item.Item, 0)
	for rows.Next() {
		var obj item.Item
		if err := rows.Scan(&obj.ID, &obj.OwnerID, &obj.Name, &obj.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, rows.Err()
}
