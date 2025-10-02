package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	ds "recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type settingsRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

type ISettingsRepository interface {
	GetSettings(userID string) (*entities.UserSettings, error)
	UpsertSettings(settings *entities.UserSettings) error
	DeleteSettings(userID string) error
}

func NewSettingsRepository(db *ds.MongoDB) ISettingsRepository {
	collection := db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("user_settings")

	return &settingsRepository{
		Collection: collection,
		Context:    db.Context,
	}
}

func (repo *settingsRepository) GetSettings(userID string) (*entities.UserSettings, error) {
	var settings entities.UserSettings
	filter := bson.M{"user_id": userID}

	err := repo.Collection.FindOne(repo.Context, filter).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default settings if not found
			return entities.DefaultUserSettings(userID), nil
		}
		return nil, fmt.Errorf("error getting settings: %v", err)
	}

	return &settings, nil
}

func (repo *settingsRepository) UpsertSettings(settings *entities.UserSettings) error {
	filter := bson.M{"user_id": settings.UserID}

	// Update the UpdatedAt timestamp
	settings.UpdatedAt = time.Now()

	// If CreatedAt is zero, set it to now
	if settings.CreatedAt.IsZero() {
		settings.CreatedAt = time.Now()
	}

	update := bson.M{"$set": settings}
	opts := options.Update().SetUpsert(true)

	_, err := repo.Collection.UpdateOne(repo.Context, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error upserting settings: %v", err)
	}

	return nil
}

func (repo *settingsRepository) DeleteSettings(userID string) error {
	filter := bson.M{"user_id": userID}

	_, err := repo.Collection.DeleteOne(repo.Context, filter)
	if err != nil {
		return fmt.Errorf("error deleting settings: %v", err)
	}

	return nil
}
