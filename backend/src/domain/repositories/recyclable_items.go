package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRecyclableItemsRepository interface {
	FindAll() (*[]entities.RecyclableItemsModel, error)
}

type recyclableItemsRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewRecyclableItemsRepository(db *ds.MongoDB) IRecyclableItemsRepository {
	return &recyclableItemsRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("recyclable_items"),
		Context:    db.Context,
	}
}

func (repo *recyclableItemsRepository) FindAll() (*[]entities.RecyclableItemsModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding recyclable items: %v", err)
	}
	defer cursor.Close(ctx)

	var recyclableItems []entities.RecyclableItemsModel
	for cursor.Next(ctx) {
		var recyclableItem entities.RecyclableItemsModel
		if err := cursor.Decode(&recyclableItem); err != nil {
			return nil, fmt.Errorf("error decoding recyclable item: %v", err)
		}
		recyclableItems = append(recyclableItems, recyclableItem)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &recyclableItems, nil
}
