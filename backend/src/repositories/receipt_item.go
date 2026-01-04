package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IReceiptItemRepository interface {
	Create(data *entities.ReceiptItem) error
	CreateMany(data []interface{}) error
	FindByReceiptID(receiptID string) (*[]entities.ReceiptItem, error)
	FindByReceiptIDs(receiptIDs []string) (*[]entities.ReceiptItem, error)
	GetPriceHistoryByName(name string) ([]entities.ReceiptItemWithDate, error)
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

func (repo *receiptItemRepository) FindByReceiptIDs(receiptIDs []string) (*[]entities.ReceiptItem, error) {
	filter := bson.M{
		"receipt_id": bson.M{
			"$in": receiptIDs,
		},
	}
	cursor, err := repo.Collection.Find(repo.Context, filter)
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

func (repo *receiptItemRepository) GetPriceHistoryByName(name string) ([]entities.ReceiptItemWithDate, error) {
	// Lookup to get created_at from receipts
	// Search by Name OR Category
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "category", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "receipts"},
			{Key: "localField", Value: "receipt_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "receipt"},
		}}},
		{{Key: "$unwind", Value: "$receipt"}},
		{{Key: "$project", Value: bson.D{
			{Key: "unit_price", Value: 1},
			{Key: "weight", Value: 1}, // Include weight for volume
			{Key: "created_at", Value: "$receipt.created_at"},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: 1}}}},
	}

	cursor, err := repo.Collection.Aggregate(repo.Context, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating price history: %v", err)
	}
	defer cursor.Close(repo.Context)

	var results []entities.ReceiptItemWithDate
	if err := cursor.All(repo.Context, &results); err != nil {
		return nil, fmt.Errorf("error decoding aggregated results: %v", err)
	}

	return results, nil
}
