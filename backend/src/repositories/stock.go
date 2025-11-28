package repositories

import (
	"context"
	"fmt"
	"os"
	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IStockRepository interface {
	UpdateStock(shopID, wasteID string, quantity float64, category, name string, pricePerKg float64) error
	GetStock(shopID, wasteID string) (*entities.Stock, error)
	GetStocksByShopID(shopID string) ([]entities.Stock, error)
	DeleteByWasteID(wasteID string) error
}

type stockRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

func NewStockRepository(db *ds.MongoDB) IStockRepository {
	return &stockRepository{
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("stocks"),
		Context:    db.Context,
	}
}

func (repo *stockRepository) UpdateStock(shopID, wasteID string, quantity float64, category, name string, pricePerKg float64) error {
	filter := bson.M{
		"shop_id":  shopID,
		"waste_id": wasteID,
	}

	update := bson.M{
		"$inc": bson.M{"quantity": quantity},
		"$set": bson.M{
			"updated_at":   time.Now(),
			"category":     category,
			"name":         name,
			"price_per_kg": pricePerKg,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := repo.Collection.UpdateOne(repo.Context, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error updating stock: %v", err)
	}

	return nil
}

func (repo *stockRepository) GetStock(shopID, wasteID string) (*entities.Stock, error) {
	filter := bson.M{
		"shop_id":  shopID,
		"waste_id": wasteID,
	}

	var stock entities.Stock
	err := repo.Collection.FindOne(repo.Context, filter).Decode(&stock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if not found
		}
		return nil, fmt.Errorf("error getting stock: %v", err)
	}

	return &stock, nil
}

func (repo *stockRepository) DeleteByWasteID(wasteID string) error {
	filter := bson.M{"waste_id": wasteID}
	_, err := repo.Collection.DeleteMany(repo.Context, filter)
	if err != nil {
		return fmt.Errorf("error deleting stock by waste_id: %v", err)
	}
	return nil
}

func (repo *stockRepository) GetStocksByShopID(shopID string) ([]entities.Stock, error) {
	filter := bson.M{"shop_id": shopID}
	cursor, err := repo.Collection.Find(repo.Context, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding stocks: %v", err)
	}
	defer cursor.Close(repo.Context)

	var stocks []entities.Stock
	if err = cursor.All(repo.Context, &stocks); err != nil {
		return nil, fmt.Errorf("error decoding stocks: %v", err)
	}

	return stocks, nil
}
