package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/mongo"
)

type IReceiptItemRepository interface {
	Create(data *entities.ReceiptItem) error
	CreateMany(data []interface{}) error
	FindByReceiptID(receiptID string) (*[]entities.ReceiptItem, error)
}

type receiptItemRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewReceiptItemRepository(db *ds.MongoDB) IReceiptItemRepository {
	return &receiptItemRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("receipt_items"),
		Context:    db.Context,
	}
}

func (repo *receiptItemRepository) Create(data *entities.ReceiptItem) error {
	_, err := repo.Collection.InsertOne(repo.Context, data)
	if err != nil {
		return fmt.Errorf("error inserting receipt item: %v", err)
	}
	return nil
}

func (repo *receiptItemRepository) CreateMany(data []interface{}) error {
	_, err := repo.Collection.InsertMany(repo.Context, data)
	if err != nil {
		return fmt.Errorf("error inserting receipt items: %v", err)
	}
	return nil
}

func (repo *receiptItemRepository) FindByReceiptID(receiptID string) (*[]entities.ReceiptItem, error) {
	cursor, err := repo.Collection.Find(repo.Context, map[string]interface{}{
		"receipt_id": receiptID,
	})
	if err != nil {
		return nil, fmt.Errorf("error finding receipt items: %v", err)
	}
	defer cursor.Close(repo.Context)

	var items []entities.ReceiptItem
	if err := cursor.All(repo.Context, &items); err != nil {
		return nil, fmt.Errorf("error decoding receipt items: %v", err)
	}

	return &items, nil
}
