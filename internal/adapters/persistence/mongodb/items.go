package mongodb

import (
	"context"
	"time"

	"example.com/go-rest-starter/internal/domain/item"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct{ collection *mongo.Collection }

func New(db *mongo.Database) *Repository { return &Repository{collection: db.Collection("items")} }
func (r *Repository) Save(ctx context.Context, obj item.Item) error {
	_, err := r.collection.InsertOne(ctx, bson.M{"_id": obj.ID, "owner_id": obj.OwnerID, "name": obj.Name, "created_at": obj.CreatedAt})
	return err
}
func (r *Repository) ListByOwner(ctx context.Context, owner string) ([]item.Item, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": owner}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	result := make([]item.Item, 0)
	for cursor.Next(ctx) {
		var doc struct {
			ID        string    `bson:"_id"`
			OwnerID   string    `bson:"owner_id"`
			Name      string    `bson:"name"`
			CreatedAt time.Time `bson:"created_at"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result = append(result, item.Item{ID: doc.ID, OwnerID: doc.OwnerID, Name: doc.Name, CreatedAt: doc.CreatedAt})
	}
	return result, cursor.Err()
}
