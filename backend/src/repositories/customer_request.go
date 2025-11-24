package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/mongo"
)

type ICustomerRequestRepository interface {
	AddCustomerRequest(userID string, body entities.CustomerRequest) error
}

type customerRequestRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewCustomerRequestRepository(db *ds.MongoDB) ICustomerRequestRepository {
	return &customerRequestRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("customer_requests"),
		Context:    db.Context,
	}
}

func (repo *customerRequestRepository) AddCustomerRequest(userID string, body entities.CustomerRequest) error {
	_, err := repo.Collection.InsertOne(repo.Context, body)
	if err != nil {
		return fmt.Errorf("error inserting customer request: %v", err)
	}
	return nil
}
