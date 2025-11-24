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

type IShopRepository interface {
	Create(data *entities.ShopModel) error
	GetByShopID(shopID string) (*entities.ShopModel, error)
	GetByUserID(userID string) (*entities.ShopModel, error)
	GetAll(page, limit int) (*[]entities.ShopModel, int64, error)
	Update(shopID string, data *entities.ShopModel) error
	Delete(shopID string) error
}

type shopRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewShopRepository(db *ds.MongoDB) IShopRepository {
	return &shopRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("shops"),
		Context:    db.Context,
	}
}

func (repo *shopRepository) Create(data *entities.ShopModel) error {
	_, err := repo.Collection.InsertOne(repo.Context, data)
	if err != nil {
		return fmt.Errorf("error inserting shop: %v", err)
	}
	return nil
}

func (repo *shopRepository) GetByShopID(shopID string) (*entities.ShopModel, error) {
	var shop entities.ShopModel
	filter := bson.M{"shop_id": shopID}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&shop); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("error finding shop: %v", err)
	}
	return &shop, nil
}

func (repo *shopRepository) GetByUserID(userID string) (*entities.ShopModel, error) {
	var shop entities.ShopModel
	filter := bson.M{"user_id": userID}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&shop); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("error finding shop by user ID: %v", err)
	}
	return &shop, nil
}

func (repo *shopRepository) GetAll(page, limit int) (*[]entities.ShopModel, int64, error) {
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
		return nil, 0, fmt.Errorf("error finding shops: %v", err)
	}
	defer cursor.Close(repo.Context)

	var shops []entities.ShopModel
	for cursor.Next(repo.Context) {
		var shop entities.ShopModel
		if err := cursor.Decode(&shop); err != nil {
			return nil, 0, fmt.Errorf("error decoding shop: %v", err)
		}
		shops = append(shops, shop)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %v", err)
	}

	// Get total count
	totalCount, err := repo.Collection.CountDocuments(repo.Context, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("error counting shops: %v", err)
	}

	return &shops, totalCount, nil
}

func (repo *shopRepository) Update(shopID string, data *entities.ShopModel) error {
	filter := bson.M{"shop_id": shopID}
	update := bson.M{"$set": data}
	if _, err := repo.Collection.UpdateOne(repo.Context, filter, update); err != nil {
		return fmt.Errorf("error updating shop: %v", err)
	}
	return nil
}

func (repo *shopRepository) Delete(shopID string) error {
	filter := bson.M{"shop_id": shopID}
	if _, err := repo.Collection.DeleteOne(repo.Context, filter); err != nil {
		return fmt.Errorf("error deleting shop: %v", err)
	}
	return nil
}
