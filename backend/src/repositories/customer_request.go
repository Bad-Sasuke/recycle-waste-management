package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICustomerRequestRepository interface {
	AddCustomerRequest(body models.CustomerRequestModel) error
	GetCustomerRequests(userID string) (*[]models.CustomerRequestModel, error)
	CheckCustomerRequestAlreadyExist(userID string) error
	DeleteAllCustomerRequest(userID string) error
	GetCustomerRequestsPublic() (*[]models.CustomerRequestModel, error)
	UpdateCustomerRequestStatus(customerRequestID string, status models.STATUS_REQUEST) error
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

func (repo *customerRequestRepository) AddCustomerRequest(body models.CustomerRequestModel) error {
	body.CreatedAt = time.Now()
	body.UpdatedAt = time.Now()
	body.Status = models.CR_PENDING
	_, err := repo.Collection.InsertOne(repo.Context, body)
	if err != nil {
		return fmt.Errorf("error inserting customer request: %v", err)
	}
	return nil
}

func (repo *customerRequestRepository) GetCustomerRequests(userID string) (*[]models.CustomerRequestModel, error) {
	cursor, err := repo.Collection.Find(repo.Context, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error getting customer requests: %v", err)
	}
	defer cursor.Close(repo.Context)
	var customerRequests []models.CustomerRequestModel
	if err := cursor.All(repo.Context, &customerRequests); err != nil {
		return nil, fmt.Errorf("error getting customer requests: %v", err)
	}
	return &customerRequests, nil
}

func (repo *customerRequestRepository) CheckCustomerRequestAlreadyExist(userID string) error {
	cursor, err := repo.Collection.Find(repo.Context, bson.M{"user_id": userID, "status": "pending"})
	if err != nil {
		return fmt.Errorf("error getting customer requests: %v", err)
	}
	defer cursor.Close(repo.Context)
	var customerRequests []models.CustomerRequestModel
	if err := cursor.All(repo.Context, &customerRequests); err != nil {
		return fmt.Errorf("error getting customer requests: %v", err)
	}
	if len(customerRequests) > 0 {
		return fmt.Errorf("customer request already exist")
	}
	return nil
}

func (repo *customerRequestRepository) DeleteAllCustomerRequest(userID string) error {
	_, err := repo.Collection.DeleteMany(repo.Context, bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("error deleting customer requests: %v", err)
	}
	return nil
}

func (repo *customerRequestRepository) GetCustomerRequestsPublic() (*[]models.CustomerRequestModel, error) {
	cursor, err := repo.Collection.Find(repo.Context, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting customer requests: %v", err)
	}
	defer cursor.Close(repo.Context)
	var customerRequests []models.CustomerRequestModel
	if err := cursor.All(repo.Context, &customerRequests); err != nil {
		return nil, fmt.Errorf("error getting customer requests: %v", err)
	}
	return &customerRequests, nil
}

func (repo *customerRequestRepository) UpdateCustomerRequestStatus(customerRequestID string, status models.STATUS_REQUEST) error {
	filter := bson.M{"customer_request_id": customerRequestID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("error updating customer request status: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("customer request not found")
	}

	return nil
}
