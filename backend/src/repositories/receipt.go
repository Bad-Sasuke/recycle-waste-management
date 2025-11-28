package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/mongo"
)

type IReceiptRepository interface {
	Create(data *entities.Receipt) error
	FindByCustomerRequestID(requestID string) (*entities.Receipt, error)
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
