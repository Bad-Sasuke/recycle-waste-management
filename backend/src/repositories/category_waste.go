package repositories

import (
	"context"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryWasteRepository interface {
	FindAll() (*[]entities.CategoryWasteModel, error)
}

type CategoryWasteRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewCategoryWasteRepository(db *ds.MongoDB) ICategoryWasteRepository {
	return &CategoryWasteRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("category_waste"),
		Context:    db.Context,
	}
}
func (r *CategoryWasteRepository) FindAll() (*[]entities.CategoryWasteModel, error) {
	var categoryWaste []entities.CategoryWasteModel
	cursor, err := r.Collection.Find(r.Context, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(r.Context, &categoryWaste); err != nil {
		return nil, err
	}
	return &categoryWaste, nil
}
