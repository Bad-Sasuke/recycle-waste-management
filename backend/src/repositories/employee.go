package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type employeeRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

type IEmployeeRepository interface {
	CreateEmployee(employee *models.Employee) error
	GetEmployeeByID(employeeID string) (*models.Employee, error)
	GetEmployeesByShopID(shopID string, page, pageSize int) ([]models.Employee, int, error)
	GetEmployeeByUsername(username string) (*models.Employee, error)
	UpdateEmployee(employeeID string, employee *models.Employee) error
	DeleteEmployee(employeeID string) error
	CountEmployeesByShopID(shopID string) (int, error)
}

func NewEmployeeRepository(db *ds.MongoDB) IEmployeeRepository {
	collection := db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("employees")

	return &employeeRepository{
		Collection: collection,
		Context:    db.Context,
	}
}

func (repo *employeeRepository) CreateEmployee(employee *models.Employee) error {
	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	if _, err := repo.Collection.InsertOne(repo.Context, employee); err != nil {
		return fmt.Errorf("error creating employee: %v", err)
	}
	return nil
}

func (repo *employeeRepository) GetEmployeeByID(employeeID string) (*models.Employee, error) {
	var employee models.Employee
	filter := bson.M{"employee_id": employeeID}

	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&employee); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("error getting employee: %v", err)
	}
	return &employee, nil
}

func (repo *employeeRepository) GetEmployeesByShopID(shopID string, page, pageSize int) ([]models.Employee, int, error) {
	filter := bson.M{"shop_id": shopID}

	// Count total documents
	total, err := repo.Collection.CountDocuments(repo.Context, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting employees: %v", err)
	}

	// Calculate skip
	skip := (page - 1) * pageSize

	// Find with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := repo.Collection.Find(repo.Context, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding employees: %v", err)
	}
	defer cursor.Close(repo.Context)

	var employees []models.Employee
	if err := cursor.All(repo.Context, &employees); err != nil {
		return nil, 0, fmt.Errorf("error decoding employees: %v", err)
	}

	return employees, int(total), nil
}

func (repo *employeeRepository) GetEmployeeByUsername(username string) (*models.Employee, error) {
	var employee models.Employee
	filter := bson.M{"username": username}

	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&employee); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting employee by username: %v", err)
	}
	return &employee, nil
}

func (repo *employeeRepository) UpdateEmployee(employeeID string, employee *models.Employee) error {
	employee.UpdatedAt = time.Now()
	filter := bson.M{"employee_id": employeeID}
	update := bson.M{"$set": employee}

	if _, err := repo.Collection.UpdateOne(repo.Context, filter, update); err != nil {
		return fmt.Errorf("error updating employee: %v", err)
	}
	return nil
}

func (repo *employeeRepository) DeleteEmployee(employeeID string) error {
	filter := bson.M{"employee_id": employeeID}

	if _, err := repo.Collection.DeleteOne(repo.Context, filter); err != nil {
		return fmt.Errorf("error deleting employee: %v", err)
	}
	return nil
}

func (repo *employeeRepository) CountEmployeesByShopID(shopID string) (int, error) {
	filter := bson.M{"shop_id": shopID}
	count, err := repo.Collection.CountDocuments(repo.Context, filter)
	if err != nil {
		return 0, fmt.Errorf("error counting employees: %v", err)
	}
	return int(count), nil
}
