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

type IRecyclableItemsRepository interface {
	FindAll() (*[]entities.RecyclableItemsModel, error)
	FindByWasteID(wasteID string) (*entities.RecyclableItemsModel, error)
	Create(data *entities.RecyclableItemsModel) error
	Delete(wasteID string) error
	Update(wasteID string, data *entities.RecyclableItemsModel) error
	FindAllPaginated(page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	FindByShopIDPaginated(shopID string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	FindByShopIDAndCategoryPaginated(shopID, category string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error)
	UpdateStock(wasteID string, quantity float64) error
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

	cursor, err := repo.Collection.Find(repo.Context, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding recyclable items: %v", err)
	}
	defer cursor.Close(repo.Context)

	var recyclableItems []entities.RecyclableItemsModel
	for cursor.Next(repo.Context) {
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

func (repo *recyclableItemsRepository) Create(data *entities.RecyclableItemsModel) error {
	_, err := repo.Collection.InsertOne(repo.Context, data)
	if err != nil {
		return fmt.Errorf("error inserting recyclable item: %v", err)
	}
	return nil
}

func (repo *recyclableItemsRepository) Delete(wasteID string) error {
	filter := bson.M{"waste_id": wasteID}
	if _, err := repo.Collection.DeleteOne(repo.Context, filter); err != nil {
		return fmt.Errorf("error deleting recyclable item: %v", err)
	}
	return nil
}

func (repo *recyclableItemsRepository) FindByWasteID(wasteID string) (*entities.RecyclableItemsModel, error) {
	var recyclableItem entities.RecyclableItemsModel
	filter := bson.M{"waste_id": wasteID}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&recyclableItem); err != nil {
		return nil, fmt.Errorf("error finding recyclable item: %v", err)
	}
	return &recyclableItem, nil
}

func (repo *recyclableItemsRepository) Update(wasteID string, data *entities.RecyclableItemsModel) error {
	filter := bson.M{"waste_id": wasteID}
	update := bson.M{"$set": bson.M{
		"shop_id":    data.ShopID,
		"name":       data.Name,
		"category":   data.Category,
		"price":      data.Price,
		"lastupdate": data.LastUpdate,
		"hours":      data.Hours,
		"url":        data.URL,
	}}
	if _, err := repo.Collection.UpdateOne(repo.Context, filter, update); err != nil {
		return fmt.Errorf("error updating recyclable item: %v", err)
	}
	return nil
}

func (repo *recyclableItemsRepository) FindAllPaginated(page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	cursor, err := repo.Collection.Find(
		repo.Context,
		bson.M{},
		&options.FindOptions{
			Skip:  &skip,
			Limit: &limit64,
		},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding recyclable items: %v", err)
	}
	defer cursor.Close(repo.Context)

	var recyclableItems []entities.RecyclableItemsModel
	for cursor.Next(repo.Context) {
		var recyclableItem entities.RecyclableItemsModel
		if err := cursor.Decode(&recyclableItem); err != nil {
			return nil, 0, fmt.Errorf("error decoding recyclable item: %v", err)
		}
		recyclableItems = append(recyclableItems, recyclableItem)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	// Get total count
	totalCount, err := repo.Collection.CountDocuments(repo.Context, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("error counting recyclable items: %v", err)
	}

	return &recyclableItems, totalCount, nil
}

func (repo *recyclableItemsRepository) FindByShopIDPaginated(shopID string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	filter := bson.M{"shop_id": shopID}
	cursor, err := repo.Collection.Find(
		repo.Context,
		filter,
		&options.FindOptions{
			Skip:  &skip,
			Limit: &limit64,
		},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding recyclable items: %v", err)
	}
	defer cursor.Close(repo.Context)

	var recyclableItems []entities.RecyclableItemsModel
	for cursor.Next(repo.Context) {
		var recyclableItem entities.RecyclableItemsModel
		if err := cursor.Decode(&recyclableItem); err != nil {
			return nil, 0, fmt.Errorf("error decoding recyclable item: %v", err)
		}
		recyclableItems = append(recyclableItems, recyclableItem)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	// Get total count
	totalCount, err := repo.Collection.CountDocuments(repo.Context, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting recyclable items: %v", err)
	}

	return &recyclableItems, totalCount, nil
}

func (repo *recyclableItemsRepository) FindByShopIDAndCategoryPaginated(shopID, category string, page, limit int) (*[]entities.RecyclableItemsModel, int64, error) {
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	filter := bson.M{"shop_id": shopID, "category": category}
	cursor, err := repo.Collection.Find(
		repo.Context,
		filter,
		&options.FindOptions{
			Skip:  &skip,
			Limit: &limit64,
		},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding recyclable items: %v", err)
	}
	defer cursor.Close(repo.Context)

	var recyclableItems []entities.RecyclableItemsModel
	for cursor.Next(repo.Context) {
		var recyclableItem entities.RecyclableItemsModel
		if err := cursor.Decode(&recyclableItem); err != nil {
			return nil, 0, fmt.Errorf("error decoding recyclable item: %v", err)
		}
		recyclableItems = append(recyclableItems, recyclableItem)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	// Get total count
	totalCount, err := repo.Collection.CountDocuments(repo.Context, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting recyclable items: %v", err)
	}

	return &recyclableItems, totalCount, nil
}

func (repo *recyclableItemsRepository) UpdateStock(wasteID string, quantity float64) error {
	filter := bson.M{"waste_id": wasteID}
	update := bson.M{"$inc": bson.M{"stock": quantity}}
	_, err := repo.Collection.UpdateOne(repo.Context, filter, update)
	if err != nil {
		return fmt.Errorf("error updating stock: %v", err)
	}
	return nil
}
