package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IReceiptRepository interface {
	Create(data *entities.Receipt) error
	FindByID(receiptID string) (*entities.Receipt, error)
	FindByCustomerRequestID(requestID string) (*entities.Receipt, error)
	FindByShopID(shopID string) ([]entities.Receipt, error)
}

type receiptRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewReceiptRepository(db *ds.MongoDB) IReceiptRepository {
	return &receiptRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("receipts"),
		Context:    db.Context,
	}
}

func (repo *receiptRepository) Create(data *entities.Receipt) error {
	_, err := repo.Collection.InsertOne(repo.Context, data)
	if err != nil {
		return fmt.Errorf("error inserting receipt: %v", err)
	}
	return nil
}

func (repo *receiptRepository) FindByID(receiptID string) (*entities.Receipt, error) {
	var receipt entities.Receipt
	err := repo.Collection.FindOne(repo.Context, map[string]interface{}{
		"_id": receiptID,
	}).Decode(&receipt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding receipt: %v", err)
	}
	return &receipt, nil
}

func (repo *receiptRepository) FindByCustomerRequestID(requestID string) (*entities.Receipt, error) {
	var receipt entities.Receipt
	err := repo.Collection.FindOne(repo.Context, map[string]interface{}{
		"customer_request_id": requestID,
	}).Decode(&receipt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding receipt: %v", err)
	}
	return &receipt, nil
}

func (repo *receiptRepository) FindByShopID(shopID string) ([]entities.Receipt, error) {
	filter := bson.M{"shop_id": shopID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}) // เรียงใหม่ไปเก่า

	cursor, err := repo.Collection.Find(repo.Context, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding receipts: %v", err)
	}
	defer cursor.Close(repo.Context)

	var receipts []entities.Receipt
	if err = cursor.All(repo.Context, &receipts); err != nil {
		return nil, fmt.Errorf("error decoding receipts: %v", err)
	}

	return receipts, nil
}
