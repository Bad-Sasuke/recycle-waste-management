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

type usersRepository struct {
	Collection *mongo.Collection
	Context    context.Context
}

type IUsersRepository interface {
	InsertNewUser(data *entities.UserDataFormat) error
	FindAll() (*[]entities.UserDataFormat, error)
	UpdateUser(userID string, data *entities.UserDataFormat) error
	UpdateUserJWT(userID string, jwt string) error
	DeleteUser(userID string) error
	GetUser(userID string) (*entities.UserDataFormat, error)
}

func NewUsersRepository(db *ds.MongoDB) IUsersRepository {
	collection := db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("users")

	return &usersRepository{
		Collection: collection,
		Context:    db.Context,
	}
}

func (repo *usersRepository) InsertNewUser(data *entities.UserDataFormat) error {
	if _, err := repo.Collection.InsertOne(repo.Context, data); err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func (repo *usersRepository) FindAll() (*[]entities.UserDataFormat, error) {

	cursor, err := repo.Collection.Find(repo.Context, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding users: %v", err)
	}
	defer cursor.Close(repo.Context)

	var userData []entities.UserDataFormat
	for cursor.Next(repo.Context) {
		var user entities.UserDataFormat
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("error decoding user: %v", err)
		}
		userData = append(userData, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &userData, nil
}

func (repo *usersRepository) UpdateUser(userID string, data *entities.UserDataFormat) error {

	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": data}

	if _, err := repo.Collection.UpdateOne(repo.Context, filter, update); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (repo *usersRepository) DeleteUser(userID string) error {

	filter := bson.M{"user_id": userID}
	if _, err := repo.Collection.DeleteOne(repo.Context, filter); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

func (repo *usersRepository) GetUser(userID string) (*entities.UserDataFormat, error) {
	var user entities.UserDataFormat
	filter := bson.M{"user_id": userID}
	if err := repo.Collection.FindOne(repo.Context, filter).Decode(&user); err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}

func (repo *usersRepository) UpdateUserJWT(userID string, jwt string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"jwt": jwt}}
	opts := options.Update().SetUpsert(true)
	if _, err := repo.Collection.UpdateOne(repo.Context, filter, update, opts); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}
